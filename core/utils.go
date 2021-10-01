package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func ClearValue(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func RegularKey(content string) (map[string]string, error) {
	re, err := regexp.Compile(`\$\{(.*?)\}`)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("regularKey error, regexp.Compile error: %s", err.Error()))
	}
	sm := re.FindAllSubmatch([]byte(content), -1)
	result := make(map[string]string)
	for _, v := range sm {
		result[string(v[1])] = string(v[0])
	}
	return result, nil
}

func ToString(v interface{}) (string, error) {
	switch v := v.(type) {
	case nil:
		return "", nil
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 32), nil
	case map[string]interface{}, []interface{}:
		nv, err := json.Marshal(v)
		if err != nil {
			return "", errors.New("marshal context value error")
		}
		return string(nv), nil
	default:
		return "", errors.Errorf("shared does not support this type '%s'", reflect.TypeOf(v))
	}
}

func CopyMap(data map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range data {
		n[k] = v
	}
	return n
}
