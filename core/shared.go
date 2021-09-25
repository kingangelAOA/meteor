package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	MissContextKeyError  = errors.New("get context value error, key not exist in context")
	ContextKeyExistError = errors.New("set context value error, key exist in context")
)

var (
	WrappedSharedPool *sync.Pool
)

func init() {
	WrappedSharedPool = &sync.Pool{}
}

type Shared struct {
	Data       map[string]interface{}
	originData map[string]interface{}
}

func NewShared(data map[string]interface{}) *Shared {
	return &Shared{
		Data:       data,
		originData: data,
	}
}

func (s *Shared) Reset() {
	s.Data = s.originData
}

func (s *Shared) Get(key string) (interface{}, error) {
	if v, ok := s.Data[key]; ok {
		return v, nil
	}
	return "", MissContextKeyError
}

func (s *Shared) UpdateBaseContent(bc *BaseContent) error {
	for k, v := range bc.Keys {
		rv, err := s.GetString(k)
		if err != nil {
			return nil
		}
		bc.Content = strings.Replace(bc.Content, v, rv, -1)
	}
	return nil
}

func (s *Shared) GetString(key string) (string, error) {
	rv, err := s.Get(key)
	if err != nil {
		return "", err
	}
	srv := ""
	switch rv := rv.(type) {
	case nil:
		srv = ""
	case string:
		srv = rv
	case int:
		srv = strconv.Itoa(rv)
	case float64:
		srv = strconv.FormatFloat(rv, 'E', -1, 32)
	case map[string]interface{}:
		v, err := json.Marshal(rv)
		if err != nil {
			return "", errors.New("marshal context value error")
		}
		srv = string(v)
	}
	return srv, nil
}

func (c *Shared) Set(key string, value interface{}) error {
	if _, ok := c.Data[key]; !ok {
		c.Data[key] = value
		return nil
	}
	return ContextKeyExistError
}

func (s *Shared) UpdateString(str string) (string, error) {
	rex := regexp.MustCompile(`\$\{(.*?)\}`)
	out := rex.FindAllStringSubmatch(str, -1)
	for _, i := range out {
		key := i[1]
		if v, ok := s.Data[key]; ok {
			fmt.Println(v)
		}
	}
	return str, nil
}

func AcquireShared() *Shared {
	v := WrappedSharedPool.Get()
	if v == nil {
		return &Shared{
			Data: map[string]interface{}{},
		}
	}
	return v.(*Shared)
}

func ReleaseShared(sr *Shared) {
	sr.Reset()
	WrappedSharedPool.Put(sr)
}
