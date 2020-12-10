package core

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	_ "net/http/pprof"
)

type WrappedTask struct {
	task   func()
	status bool
}

type WrappedPool struct {
	task       chan WrappedTask
	p          *ants.Pool
	readyClose chan int
}

func NewPool(num int) (*WrappedPool, error) {
	if p, err := ants.NewPool(num, ants.WithPreAlloc(true)); err != nil {
		return nil, err
	} else {
		return &WrappedPool{
			p:          p,
			task:       make(chan WrappedTask),
			readyClose: make(chan int),
		}, nil
	}
}

func (rp *WrappedPool) AddTask(wt WrappedTask) {
	rp.task <- wt
}

func (rp *WrappedPool) Run() {
	go func(rp *WrappedPool) {
		for {
			wt := <-rp.task
			if !wt.status {
				break
			}
			if err := rp.p.Submit(wt.task); err != nil {
				fmt.Println("******************************", err.Error())
			}
		}
		rp.p.Release()
	}(rp)
}
