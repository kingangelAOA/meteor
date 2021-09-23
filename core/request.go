package core

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

type RequestAuthentication map[string]string

type RequestConfig struct {
	Url            string
	Method         string
	Body           RequestBody
	Parameters     []RequestParameter
	Headers        []RequestHeader
	Authentication RequestAuthentication
}

func (rc *RequestConfig) GetRequest() *fasthttp.Request {
	// request := fasthttp.AcquireRequest()
	rc.updateContentType()
	return nil

}

func (rc *RequestConfig) updateContentType() {
	contentTypeHeaderValue := rc.Body.MimeType
	if rc.Body.MimeType == CONTENT_TYPE_GRAPHQL {
		contentTypeHeaderValue = CONTENT_TYPE_JSON
		rc.Method = METHOD_POST
	}
	hasBody := false
	if rc.Body.MimeType != "" {
		hasBody = true
	}
	rc.Headers = delContentType(rc.Headers)
	if hasBody {
		rc.Headers = append(rc.Headers, RequestHeader{Name: "Content-Type", Value: contentTypeHeaderValue})
	}
}

func (rc *RequestConfig) getBody() ([]byte, error) {
	if rc.Body.MimeType == CONTENT_TYPE_FORM_URLENCODED {
		if p, err := json.Marshal(rc.Body.Params); err != nil {
			return nil, &HttpError{Err: FormURLEncodeParamsMarshalToByteError}
		} else {
			return p, nil
		}
	} else if rc.Body.MimeType == CONTENT_TYPE_FORM_DATA {
		// bodyBuffer := &bytes.Buffer{}
		// bodyWriter := multipart.NewWriter(bodyBuffer)
		// defer func() {
		// 	contentType := bodyWriter.FormDataContentType()
		// 	headers.SetContentType(contentType)
		// }()
	}
	return nil, nil
}

type RequestBody struct {
	MimeType string
	Text     string
	File     string
	Params   []RequestBodyParameter
}

func (rb *RequestBody) getMultipartFromBody(headers []RequestHeader) map[string]string {
	// bodyBuffer := &bytes.Buffer{}
	// bodyWriter := multipart.NewWriter(bodyBuffer)

	// for _, param := range rbps {
	// 	if param.Type == MULTIPART_TYPE_FILE {

	// 	}
	// }
	return nil
}

type RequestParameter struct {
	Name  string
	Value string
}

type RequestBodyParameter struct {
	Name      string
	Value     string
	Multiline string
	File      string
	Type      string
}

type RequestHeader struct {
	Name  string
	Value string
}

func NewAuth(t string, oldAuth RequestAuthentication) RequestAuthentication {
	switch t {
	case AUTH_NONE:
		return map[string]string{}
	case AUTH_BASIC:
		return getBasicAuth(oldAuth)
	case AUTH_DIGEST:
		return getBasicAuth(oldAuth)
	case AUTH_NTLM:
		return getBasicAuth(oldAuth)
	default:
		return map[string]string{}
	}
}

func getBasicAuth(oldAuth RequestAuthentication) map[string]string {
	return map[string]string{
		"username": oldAuth["username"],
		"password": oldAuth["password"],
	}
}

func delContentType(headers []RequestHeader) []RequestHeader {
	for i := 0; i < len(headers); i++ {
		if strings.ToLower(headers[i].Name) == "content-type" {
			headers = append(headers[:i], headers[i+1:]...)
			i--
		}
	}
	return headers
}
