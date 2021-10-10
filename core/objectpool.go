package core

import (
	"context"
	"io"
	"net/http"
	"sync"
)

var (
	RequestPool       sync.Pool
	HttpMessagePool   sync.Pool
	TengoVMPool       = make(map[string]sync.Pool)
	ScriptMessagePool sync.Pool
	SharedPool        sync.Pool
	SharedMapPool     sync.Pool
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

func AcquireHttpMessage(hr *http.Request, data map[string]interface{}) *HttpMessage {
	v := HttpMessagePool.Get()
	if v == nil {
		return NewHttpMessage(hr, data)
	}
	return v.(*HttpMessage)
}

func ReleaseHttpMessage(hm *HttpMessage) {
	hm.Reset()
	RequestPool.Put(hm)
}

func AcquireTengoVM(key, code string) *TengoVM {
	if pool, ok := TengoVMPool[key]; ok {
		v := pool.Get()
		if v == nil {
			return &TengoVM{
				key:    key,
				script: newTengoVM(code),
			}
		}
		return v.(*TengoVM)
	} else {
		TengoVMPool[key] = sync.Pool{}
		return &TengoVM{
			key:    key,
			script: newTengoVM(code),
		}
	}
}

func ReleaseTengoVM(tvm *TengoVM) {
	tvm.Reset()
	if v, ok := TengoVMPool[tvm.key]; ok {
		v.Put(tvm)
	} else {
		pool := sync.Pool{}
		pool.Put(tvm)
		TengoVMPool[tvm.key] = pool
	}
}

func AcquireScriptMessage(name, t string, data map[string]interface{}) *ScriptMessage {
	v := ScriptMessagePool.Get()
	if v == nil {
		return NewScriptMessage(name, t, data)
	}
	tm := v.(*ScriptMessage)
	tm.Key = name
	tm.Type = t
	tm.SetData(data)
	return tm
}

func ReleaseScriptMessage(tm *ScriptMessage) {
	tm.Reset()
	ScriptMessagePool.Put(tm)
}

// func AcquireShared() *Shared {
//  v := SharedPool.Get()
//  if v == nil {
//   ns := &Shared{
//    Data: make(map[string]interface{}),
//   }
//   return ns
//  }
//  s := v.(*Shared)
//  return s
// }

// func ReleaseShared(s *Shared) {
//  ClearMap(s.Data)
//  SharedPool.Put(s)
// }
