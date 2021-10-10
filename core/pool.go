package core

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

type WrappedPool struct {
	p      *ants.Pool
	ctx    context.Context
	ErrMsg chan string
	task   func()
}

func NewPool(num int, ctx context.Context, task func()) (*WrappedPool, error) {
	wp := &WrappedPool{
		ctx:    ctx,
		ErrMsg: make(chan string, ErrMsgThreshold),
		task:   task,
	}
	if p, err := ants.NewPool(num, ants.WithPanicHandler(func(p interface{}) {
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

func (rp *WrappedPool) RunByLimit(l Limiter) {
	go func(l Limiter) {
		for {
			err := l.Get()
			if err != nil {
				rp.p.Release()
				break
			}
			rp.baseRun()
		}
	}(l)

}

func (rp *WrappedPool) Run() {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				rp.p.Release()
				return
			default:
				rp.baseRun()
			}
		}
	}(rp.ctx)
}

func (wp *WrappedPool) baseRun() {
	if err := wp.p.Submit(wp.task); err != nil {
		wp.PutErrMsg(err.Error())
	}
}
