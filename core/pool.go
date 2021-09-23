package core

import (
	"context"
	"fmt"

	"github.com/panjf2000/ants/v2"
)

type WrappedPool struct {
	p   *ants.PoolWithFunc
	ctx context.Context
}

func NewPool(num int, ctx context.Context, task func(wc interface{})) (*WrappedPool, error) {
	if p, err := ants.NewPoolWithFunc(num, func(i interface{}) {
		task(i)
	}, ants.WithPreAlloc(true)); err != nil {
		return nil, err
	} else {
		return &WrappedPool{
			p:   p,
			ctx: ctx,
		}, nil
	}
}

func (rp *WrappedPool) SetNum(n int) {
	rp.p.Tune(n)
}

func (rp *WrappedPool) RunByLimit(l Limiter, data interface{}) {
	go func(l Limiter, data interface{}) {
		for {
			err := l.Get()
			if err != nil {
				rp.p.Release()
				break
			}
			rp.baseRun(data)
		}
	}(l, data)

}

func (rp *WrappedPool) Run(data interface{}) {
	go func(ctx context.Context, data interface{}) {
		for {
			select {
			case <-ctx.Done():
				rp.p.Release()
				return
			default:
				rp.baseRun(data)
			}
		}
	}(rp.ctx, data)
}

func (rp *WrappedPool) baseRun(data interface{}) {
	if err := rp.p.Invoke(data); err != nil {
		fmt.Println("******************************", err.Error())
	}
}
