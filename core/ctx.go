package core

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

var (
	MissContextKeyError  = errors.New("get context value error, key not exist in context")
	ContextKeyExistError = errors.New("set context value error, key exist in context")
)

var (
	WrappedContextPool sync.Pool
)

const (
	LeftFlag  = "${"
	RightFlag = "}"
)

type WrappedContext struct {
	Data map[string]interface{}
}

func (wc *WrappedContext) Reset() {
	wc.Data = map[string]interface{}{}
}

func (wc *WrappedContext) Get(key string) (interface{}, error) {
	if v, ok := wc.Data[key]; ok {
		return v, nil
	}
	return "", MissContextKeyError
}

func (c *WrappedContext) Set(key, value string) error {
	if _, ok := c.Data[key]; !ok {
		c.Data[key] = value
		return nil
	}
	return ContextKeyExistError
}

func (c *WrappedContext) UpdateString(str string) (string, error) {
	rex := regexp.MustCompile(`\$\{(.*?)\}`)
	out := rex.FindAllStringSubmatch(str, -1)
	for _, i := range out {
		key := i[1]
		if v, ok := c.Data[key]; ok {
			fmt.Println(v)
		}
	}
	return str, nil
}

func AcquireWrappedContext() *WrappedContext {
	v := WrappedContextPool.Get()
	if v == nil {
		return &WrappedContext{
			Data: map[string]interface{}{},
		}
	}
	return v.(*WrappedContext)
}

func ReleaseWrappedContext(sr *WrappedContext) {
	sr.Reset()
	WrappedContextPool.Put(sr)
}
