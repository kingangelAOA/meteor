package core

import (
	"encoding/json"
	"strconv"
	"testing"
)

func BenchmarkHttp(b *testing.B) {
	var a float64

	a = 1.324234234234
	for n := 0; n < b.N; n++ {
		strconv.FormatFloat(a, 'E', -1, 32)
	}
}

func BenchmarkHttp1(b *testing.B) {
	var a interface{}

	a = 1.324234234234
	for n := 0; n < b.N; n++ {
		json.Marshal(a)
	}
}
