package core

import (
	"context"
	"fmt"
	"math/rand"
	_ "net/http/pprof"
	"testing"
	"time"
)

func test() {
	t := rand.Intn(1000)
	local := time.Now()
	time.Sleep(time.Duration(t) * time.Millisecond)

	fmt.Println("Hello World!*********", t, local)
}

func TestNewPool(t *testing.T) {
	begin := time.Now()
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	wc := &WrappedContext{
		Data: map[string]interface{}{},
	}
	defer cancel()
	l := NewLimiter(LimitConstantMode, 5, ctx)
	p, _ := NewPool(1, ctx, task)
	p.RunByLimit(l, wc)
	fmt.Println("run time", time.Since(begin))
	time.Sleep(15 * time.Second)
}

func task(i interface{}) {
	time.Sleep(2000 * time.Millisecond)
	// fmt.Println("this is task")
}
