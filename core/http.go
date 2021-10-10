package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
)

var (
	ErrorNewHttpServerPool = errors.New("new http pool error")
)

func init() {
	// RequestPool = sync.Pool{}
}

const (
	TypeMultiPartFormFile = "TypeMultiPartFormFile"
	TypeMultiPartFormText = "TypeMultiPartFormText"
	TypeFormURLEncode     = "TypeFormURLEncode"
	TypeJson              = "TypeJson"
	TypeBinaryFile        = "TypeBinaryFile"
	TypeNoBody            = "TypeNoBody"

	InitialHttpServicePool = 1000
	RequestThreshold       = 1000
	DefaultTimeout         = 10000
)

type HttpService struct {
	BaseService
	client *Http
}

func NewHttpService(timeout, size int, ctx context.Context) (*HttpService, error) {
	hs := &HttpService{
		client: NewHttpClient(timeout),
	}
	p, err := ants.NewPool(size, ants.WithPanicHandler(func(p interface{}) {
		hs.PutErrMsg(fmt.Sprintf("worker exits from a panic: %v\n", p))
	}), ants.WithExpiryDuration(ExpiryDuration*time.Minute))
	if err != nil {
		return nil, ErrorNewHttpServerPool
	}
	bs := NewBaseService(p, ctx)
	hs.BaseService = *bs
	return hs, nil
}

func NewHttpClient(timeout int) *Http {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &Http{
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Duration(timeout) * time.Millisecond,
		},
	}
}

func (hs *HttpService) Run() {
	go hs.requestQueue()
}

func (hs *HttpService) requestQueue() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			hs.Code = 500
			hs.PutErrMsg(err.Error())
		}
	}()
	for {
		select {
		case <-hs.ctx.Done():
			return
		case hm := <-hs.ms:
			if err := hs.p.Submit(func() {
				hs.client.Execute(hm)
			}); err != nil {
				hm.SetErr(err.Error())
			}
		}
	}
}

type Http struct {
	client *http.Client
}

func (h *Http) Execute(m Message) {
	defer m.TimeCost()
	hm := m.(*HttpMessage)
	res, err := h.client.Do(hm.Request)
	if err != nil {
		hm.SetErr(err.Error())
		return
	}
	defer ReleaseRequest(hm.Request)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hm.SetErr(err.Error())
		return
	}
	hm.Body = string(body)
	hm.ResHeader = res.Header
	hm.ResCookies = parseCookies(res.Cookies())
	hm.SetOk(true)
}

type HttpMessage struct {
	BaseMessage
	Key        string
	Request    *http.Request
	ResHeader  http.Header
	ResCookies map[string]string
	Body       string
}

func NewHttpMessage(hr *http.Request, data map[string]interface{}) *HttpMessage {
	return &HttpMessage{
		BaseMessage: NewBaseMessage(data),
		Request:     hr,
		ResHeader:   make(http.Header),
		ResCookies:  map[string]string{},
	}
}

func (hm *HttpMessage) GetKey() string {
	return hm.Key
}

func (hm *HttpMessage) Reset() {
	hm.Body = ""
	hm.ResHeader = make(http.Header)
	hm.ResCookies = map[string]string{}
	hm.ErrMsg = ""
	hm.Prints = ""
	hm.BaseMessage.reset()
}

func (hm *HttpMessage) GetCookie(key string) (string, error) {
	if v, ok := hm.ResCookies[key]; ok {
		return v, nil
	} else {
		return "", errors.New(fmt.Sprintf("cookie '%s' is not exist", key))
	}
}

func parseCookies(cookies []*http.Cookie) map[string]string {
	cks := make(map[string]string)
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		cks[cookie.Name] = cookie.Value
	}
	return cks
}

type WrappedRequest struct {
	Name   string
	Method string
	Host   string
	Port   string
	Path   *BaseContent
	Body   Body
	Header map[string]*BaseContent
	Query  map[string]*BaseContent
	Ctx    context.Context
}

func (wr *WrappedRequest) Reset() {
	wr.Path.Reset()
	wr.Body.Reset()
	for _, v := range wr.Header {
		v.Reset()
	}
	for _, v := range wr.Query {
		v.Reset()
	}
}

func (wr WrappedRequest) GetRequest(s *Shared) (*http.Request, error) {
	body, contentType, err := wr.Body.GetContent(s)
	if err != nil {
		return nil, err
	}
	r, err := AcquireRequest(wr.Method, wr.getUrl(s), body, wr.Ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range wr.Header {
		s.UpdateBaseContent(v)
		r.Header.Add(k, v.Content)
	}
	r.Header.Set("ContentType", contentType)
	q := r.URL.Query()
	for k, v := range wr.Query {
		s.UpdateBaseContent(v)
		q.Add(k, v.Content)
	}
	return r, nil
}

func (wr *WrappedRequest) getUrl(s *Shared) string {
	s.UpdateBaseContent(wr.Path)
	if wr.Port != "" {
		return fmt.Sprintf("%s:%s%s", wr.Host, wr.Port, wr.Path)
	}
	return fmt.Sprintf("%s%s", wr.Host, wr.Path)
}

type BaseContent struct {
	Content string
	Keys    map[string]string
}

func (bc *BaseContent) Reset() {
	bc.Content = ""
}

func (bc *BaseContent) UpdateContent(s *Shared) error {
	for k, v := range bc.Keys {
		rv, err := s.GetString(k)
		if err != nil {
			return err
		}
		bc.Content = strings.Replace(bc.Content, v, rv, -1)
	}
	return nil
}

func NewBaseContent(content string) (*BaseContent, error) {
	keys, err := RegularKey(content)
	if err != nil {
		return nil, err
	}
	return &BaseContent{
		Content: content,
		Keys:    keys,
	}, nil
}

type Body interface {
	GetContent(*Shared) (io.ReadCloser, string, error)
	Reset()
}

type MultipartForm struct {
	Content map[string]*MultipartValue
}

type MultipartValue struct {
	Type  string
	Value interface{}
	Keys  map[string]string
}

func (mv *MultipartValue) reset() {
	ClearValue(mv.Value)
}

func NewMultipartValue(t string, value interface{}) *MultipartValue {
	return &MultipartValue{
		Type:  t,
		Value: value,
		Keys:  make(map[string]string),
	}
}

func NewMultipartForm(content map[string]*MultipartValue) (*MultipartForm, error) {
	for _, v := range content {
		if v.Type == TypeMultiPartFormText {
			keys := map[string]string{}
			subKeys, err := RegularKey(v.Value.(string))
			if err != nil {
				return nil, err
			}
			for sk, sv := range subKeys {
				keys[sk] = sv
			}
			v.Keys = keys
		}
	}
	return &MultipartForm{
		Content: content,
	}, nil
}

func (mf *MultipartForm) GetContent(s *Shared) (io.Reader, string, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	contentType := bodyWriter.FormDataContentType()
	for k, mv := range mf.Content {
		if mv.Type == TypeMultiPartFormText {
			mvv, ok := mv.Value.(string)
			if !ok {
				return nil, "", errors.New("multipartFormText value must be string")
			}
			for k, v := range mv.Keys {
				rv, err := s.GetString(k)
				if err != nil {
					return nil, "", err
				}
				mvv = strings.Replace(mvv, v, rv, -1)
			}
			bodyWriter.WriteField(k, mvv)
		} else if mv.Type == TypeMultiPartFormFile {
			file, ok := mv.Value.(os.File)
			defer file.Close()
			if !ok {
				return nil, "", errors.New("multipartFormFile value must be File")
			}
			wf, err := bodyWriter.CreateFormFile(k, file.Name())
			if err != nil {
				return nil, "", errors.New("create form file error")
			}
			b := []byte{}
			if _, err := file.Read(b); err != nil {
				return nil, "", errors.New("read file error")
			}
			if _, err := wf.Write(b); err != nil {
				return nil, "", errors.New("write multipart form file error")
			}
		} else {
			return nil, "", errors.Errorf("does not support this type '%s'", mv.Type)
		}
	}
	return bodyBuffer, contentType, nil
}

func (mf *MultipartForm) Reset() {
	for _, v := range mf.Content {
		v.reset()
	}
}

type FormURLEncode struct {
	Content map[string]*BaseContent
}

func (fue *FormURLEncode) GetContent(s *Shared) (io.Reader, string, error) {
	DataUrlVal := url.Values{}
	for k, bc := range fue.Content {
		s.UpdateBaseContent(bc)
		DataUrlVal.Add(k, bc.Content)
	}
	return strings.NewReader(DataUrlVal.Encode()), "application/x-www-form-urlencoded", nil
}

type JsonBody struct {
	BaseContent
}

func NewJsonBody(content string) (*JsonBody, error) {
	bc, err := NewBaseContent(content)
	if err != nil {
		return nil, err
	}
	return &JsonBody{
		BaseContent: *bc,
	}, nil
}

func (jb *JsonBody) GetContent(s *Shared) (io.Reader, string, error) {
	s.UpdateBaseContent(&jb.BaseContent)
	return strings.NewReader(jb.Content), "application/json", nil
}

type Binary struct {
	Data os.File
}

func (b *Binary) GetContent() (io.Reader, string, error) {
	bs := []byte{}
	_, err := b.Data.Read(bs)
	if err != nil {
		return nil, "", errors.New("read binary file error")
	}
	return bytes.NewReader(bs), "application/json", nil
}

type Auth interface {
}
