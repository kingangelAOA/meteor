package engine_test

import (
	"testing"
)

type Test struct {
	data int
}

func (t Test) set(a int) int {
	t.data = a
	return t.data
}
func (t *Test) NewSetet(a int) {
	t.data = a
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func ReverseLinked() {
	
}

func TestName(t *testing.T) {
	
}
