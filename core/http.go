package core

import "net/http"
import "github.com/valyala/fasthttp"

type WrappedHttp struct {

}

func (wh *WrappedHttp) Request(config SamplerConfig) Result {
	client := http.Client{}
	client.Do()
	fasthttp.Do()
}

func (wh *WrappedHttp) Init() {

}
