package core

import (
	"fmt"
	"testing"
)

type A struct {
	content string
}

type B struct {
	a map[string]A
}

func newB(c string) *B {
	return &B{
		a: map[string]A{
			"test": A{content: c},
		},
	}
}

func change(b *B) {
	a := b.a["test"]
	a.content = "sdfsdfsdfsdf"
}

func TestStruct(t *testing.T) {
	b := newB("cccccc")
	fmt.Println(b.a["test"].content)
	change(b)
	fmt.Println(b.a["test"].content)
}
