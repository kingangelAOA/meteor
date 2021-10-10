package core

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	MissContextKeyError  = errors.New("get context value error, key not exist in context")
	ContextKeyExistError = errors.New("set context value error, key exist in context")
)

type Shared struct {
	Data map[string]interface{}
}

func NewShared(data map[string]interface{}) *Shared {
	return &Shared{
		Data: data,
	}
}

func (s *Shared) CopyShared() *Shared {
	ns := &Shared{
		Data: make(map[string]interface{}),
	}
	ns.SetData(s.Data)
	return ns
}

func (s *Shared) Reset() {
	s.Data = make(map[string]interface{})
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
	return ToString(rv)
}

func (s *Shared) Set(key string, value interface{}) {
	s.Data[key] = value
}

func (s *Shared) SetData(data map[string]interface{}) {
	for k, v := range data {
		s.Set(k, v)
	}
}

// func (s *Shared) UpdateString(str string) (string, error) {
// 	rex := regexp.MustCompile(`\$\{(.*?)\}`)
// 	out := rex.FindAllStringSubmatch(str, -1)
// 	for _, i := range out {
// 		key := i[1]
// 		if v, ok := s.Data[key]; ok {
// 			fmt.Println(v)
// 		}
// 	}
// 	return str, nil
// }
