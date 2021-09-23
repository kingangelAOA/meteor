package core

import (
	"context"

	"golang.org/x/time/rate"
)

const (
	LimitConstantMode = "LIMIT_CONSTANT_MODE"
)

type LimiterType struct {
	T string
}

type Limiter struct {
	l           *rate.Limiter
	ctx         context.Context
	LimiterType LimiterType
}

func NewLimiter(t string, per int, ctx context.Context) Limiter {
	return Limiter{
		l: rate.NewLimiter(rate.Limit(per), 1),
		LimiterType: LimiterType{
			T: t,
		},
		ctx: ctx,
	}
}

func (limiter *Limiter) Get() error {
	return limiter.l.Wait(limiter.ctx)
}
