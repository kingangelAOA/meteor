package engine

import (
	"context"

	"golang.org/x/time/rate"
)

type Limiter struct {
	l *rate.Limiter
}

func NewLimiter(per int) Limiter {
	return Limiter{
		l: rate.NewLimiter(rate.Limit(per), getBursts(per)),
	}
}

func (limiter *Limiter) Get(ctx context.Context) error {
	// atomic.AddInt64(&limiter.num, 1)
	// num := atomic.LoadInt64(&limiter.num)
	// if num%10000 == 0 {
	// 	now := time.Now().UnixMicro()
	// 	fmt.Println(now)
	// 	fmt.Println(num)
	// }
	return limiter.l.Wait(ctx)
}

func getBursts(per int) int {
	b := per / 10
	if b == 0 {
		b = 1
	}
	return b
}
