package core

import (
	"context"
	"io"
	"net/http"
	"sync"
)

var (
	RequestPool      sync.Pool
	HttpMessagePool  sync.Pool
	TengoVMPool      = make(map[string]sync.Pool)
	TengoMessagePool sync.Pool
	SharedPool       sync.Pool
)

func AcquireRequest(method, url string, body io.Reader, ctx context.Context) (*http.Request, error) {
	v := RequestPool.Get()
	if v == nil {
		return http.NewRequestWithContext(ctx, method, url, body)
	}
	return v.(*http.Request), nil
}

func ReleaseRequest(hr *http.Request) {
	RequestPool.Put(hr)
}

func AcquireHttpMessage(hr *http.Request, s *Shared) *HttpMessage {
	v := HttpMessagePool.Get()
	if v == nil {
		return NewHttpMessage(hr, s)
	}
	return v.(*HttpMessage)
}

func ReleaseHttpMessage(hm *HttpMessage) *Shared {
	s := hm.Reset()
	RequestPool.Put(hm)
	return s
}

func AcquireTengoVM(name, code string) *TengoVM {
	if pool, ok := TengoVMPool[name]; ok {
		v := pool.Get()
		if v == nil {
			return &TengoVM{
				name:   name,
				script: newTengoVM(code),
			}
		}
		return v.(*TengoVM)
	} else {
		TengoVMPool[name] = sync.Pool{}
		return &TengoVM{
			name:   name,
			script: newTengoVM(code),
		}
	}
}

func ReleaseTengoVM(tvm *TengoVM) {
	tvm.Reset()
	if v, ok := TengoVMPool[tvm.name]; ok {
		v.Put(tvm)
	} else {
		pool := sync.Pool{}
		pool.Put(tvm)
		TengoVMPool[tvm.name] = pool
	}
}

func AcquireTengoMessage(name string, s *Shared) *TengoMessage {
	v := TengoMessagePool.Get()
	if v == nil {
		return NewTengoMessage(name, s)
	}
	tm := v.(*TengoMessage)
	tm.Name = name
	tm.s = s
	return tm
}

func ReleaseTengoMessage(tm *TengoMessage) *Shared {
	s := tm.Reset()
	TengoMessagePool.Put(tm)
	return s
}

func AcquireShared(data map[string]interface{}) *Shared {
	v := SharedPool.Get()
	if v == nil {
		return &Shared{
			Data: data,
		}
	}
	s := v.(*Shared)
	s.Data = data
	return s
}

func ReleaseShared(s *Shared) {
	s.Reset()
	SharedPool.Put(s)
}
