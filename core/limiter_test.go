package core

import (
	"context"
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	limiter := NewLimiter(LimitConstantMode, 1000, ctx)
	// defer cancel()
	index := 0
	for {
		err := limiter.Get()
		if err != nil {
			break
		}
		index += 1
		if index%1000 == 0 {
			fmt.Println(index)
		}
		if index == 3000 {
			// ctx.Done()
			cancel()
		}
	}
}
