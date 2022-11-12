package common_test

import (
	"common"
	"fmt"
	"testing"
)

func TestInsertByIndex(t *testing.T) {
	a := []int{1,2,3,4,6}
	r := common.InsertByIndex(11111, 2, a)
	fmt.Println(r)
}