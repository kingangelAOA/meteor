package core

import (
	"fmt"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	b := time.Now()
	time.Sleep(1 * time.Millisecond)
	e := time.Since(b).Nanoseconds()
	fmt.Println(e, b.UnixMilli(), b.Unix())
}
