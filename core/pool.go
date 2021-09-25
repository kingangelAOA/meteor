package core

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

type WrappedPool struct {
	p      *ants.PoolWithFunc
	ctx    context.Context
	ErrMsg chan string
}

func NewPool(num int, ctx context.Context, task func(wc interface{})) (*WrappedPool, error) {
	wp := &WrappedPool{
		ctx:    ctx,
		ErrMsg: make(chan string, ErrMsgThreshold),
	}
	if p, err := ants.NewPoolWithFunc(num, func(i interface{}) {
		task(i)
	}, ants.WithPanicHandler(func(p interface{}) {
		wp.PutErrMsg(fmt.Sprintf("worker exits from a panic: %v\n", p))
	}), ants.WithExpiryDuration(ExpiryDuration*time.Minute)); err != nil {
		return nil, err
	} else {
		wp.p = p
		return wp, nil
	}
}

func (wp *WrappedPool) PutErrMsg(m string) {
	if len(wp.ErrMsg) > ErrMsgThreshold {
		<-wp.ErrMsg
	}
	wp.ErrMsg <- m
}

func (wp *WrappedPool) SetNum(n int) {
	wp.p.Tune(n)
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

func (wp *WrappedPool) baseRun(data interface{}) {
	if err := wp.p.Invoke(data); err != nil {
		wp.PutErrMsg(err.Error())
	}
}
