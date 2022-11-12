package engine

import (
	"common"
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

type WrappedPool struct {
	p      *ants.Pool
	ErrMsg chan string
	task   func()
}

func NewPool(num int, task func()) (*WrappedPool, error) {
	wp := &WrappedPool{
		ErrMsg: make(chan string, common.ErrMsgThreshold),
		task:   task,
	}
	if p, err := ants.NewPool(num, ants.WithPanicHandler(func(p interface{}) {
		wp.PutErrMsg(fmt.Sprintf("worker exits from a panic: %v\n", p))
	}), ants.WithExpiryDuration(common.ExpiryDuration*time.Second)); err != nil {
		return nil, err
	} else {
		wp.p = p
		return wp, nil
	}
}

func (wp *WrappedPool) PutErrMsg(m string) {
	if !wp.p.IsClosed() {
		if len(wp.ErrMsg) > common.ErrMsgThreshold {
			<-wp.ErrMsg
		}
		wp.ErrMsg <- m
	}
}

func (wp *WrappedPool) close() {
	close(wp.ErrMsg)
}

func (wp *WrappedPool) SetNum(n int) {
	wp.p.Tune(n)
}

func (wp *WrappedPool) RunByLimit(ctx context.Context, l Limiter) {
	go func(ctx context.Context, l Limiter) {
		for {
			select {
			case <-ctx.Done():
				if err := wp.p.ReleaseTimeout(1 * time.Millisecond); err != nil {
					fmt.Println(err.Error())
				}
				wp.close()
				return
			default:
				err := l.Get(ctx)
				if err == nil {
					wp.baseRun()
				}
			}
		}
	}(ctx, l)
}

func (wp *WrappedPool) Run(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				if err := wp.p.ReleaseTimeout(1 * time.Millisecond); err != nil {
					fmt.Println(err.Error())
				}
				wp.close()
				return
			default:
				wp.baseRun()
			}
		}
	}(ctx)
}

func (wp *WrappedPool) baseRun() {
	if err := wp.p.Submit(wp.task); err != nil {
		fmt.Println(err.Error())
	}
}
