package core

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	TypeMultiPartForm = "TypeMultiPartForm"
	TypeFormURLEncode = "TypeFormURLEncode"
	TypeText          = "TypeText"
	TypeBinaryFile    = "TypeBinaryFile"
	TypeNoBody        = "TypeNoBody"
)

type WrappedHttp struct {
	Method      string
	Host        string
	Port        string
	Path        string
	Body        Body
	Header      map[string]string
	Query       map[string]string
	Timeout     int
	client      *http.Client
	context     *WrappedContext
	requestPool sync.Pool
}

type HttpResult struct {
}

func (wh *WrappedHttp) acquireRequest() (*http.Request, error) {
	v := wh.requestPool.Get()
	if v == nil {
		return http.NewRequest(wh.Method, wh.getUrl(), nil)
	}
	return v.(*http.Request), nil
}

func (wh *WrappedHttp) releaseRequest(req *http.Request) {
	req.Body = nil
}

func (wh *WrappedHttp) Request(config SamplerConfig) Result {
	return HttpResult{}
}

func (wh *WrappedHttp) UpdateContext(context *WrappedContext) {
	wh.context = context
}

// func (wh *WrappedHttp) newRequest() (*http.Request, error) {

//  if wh.Method == http.MethodGet {

//  }

//  req, err := http.NewRequest(wh.Method, wh.getUrl(), nil)
//  if err != nil {
//   return nil, err
//  }
// }

func (wh *WrappedHttp) getUrl() string {
	if wh.Port != "" {
		return fmt.Sprintf("%s:%s%s", wh.Host, wh.Port, wh.Path)
	}
	return fmt.Sprintf("%s%s", wh.Host, wh.Path)
}

func (wh *WrappedHttp) Init() (err error) {
	wh.client = &http.Client{
		Timeout: time.Millisecond * time.Duration(wh.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return
}

func (wh *WrappedHttp) getBody() Body {
	return nil
}

type Body interface {
	GetContent(ctx *WrappedContext) (io.Reader, error)
	GetType() string
}

type BodyExtra struct {
	bodyType string
	keys     map[string]string
}

func (be *BodyExtra) GetType() string {
	return be.bodyType
}

func (be *BodyExtra) get() string {
	return be.bodyType
}

type MultipartForm struct {
	Content map[string]MultipartValue
	BodyExtra
}

type MultipartValue struct {
	Type  string
	Value string
}

type FormURLEncode struct {
	Content map[string]string
	BodyExtra
}

type Text struct {
	Content string
	BodyExtra
}

func (t *Text) GetContent(ctx *WrappedContext) (io.Reader, error) {
	for k, v := range t.keys {
		rv, err := ctx.Get(k)
		if err != nil {
			return nil, err
		}
		t.Content = strings.Replace(t.Content, v, rv.(string), -1)
	}
	return strings.NewReader(t.Content), nil
}

type Binary []byte

type Auth interface {
}

type Query map[string]string

type Header map[string]string
