package common

import (
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetLocalZone() *time.Location {
	local, _ := time.LoadLocation("Local")
	return local
}

func MergeBsonM(left, right bson.M) bson.M {
	for k, v := range right {
		left[k] = v
	}
	return left
}

func GetMapKeys(m map[string]interface{}) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsClosed[T any](ch <-chan T) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}

func InsertByIndex[T any](e T, i int, arr []T)  []T  {
	arr = append(arr[:i+1], arr[i:]...)
	arr[i] = e
	return arr
}