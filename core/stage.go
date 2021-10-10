package core

import (
	"context"
	"errors"
)

const (
	SingleGoroutine   = "SingleGoroutine"
	MultipleGoroutine = "MultipleGoroutine"
)

type Stage interface {
	Run(p Shared)
}

type SingleCoroutineStage struct {
	s         *Shared
	beginNode Node
	ctx       context.Context
}

func ConnectSingleCoroutineNode(nodes []Node, ctx context.Context, s *Shared, collect bool) (*SingleCoroutineStage, error) {
	size := len(nodes)
	if size == 0 {
		return nil, errors.New("must be at least one node")
	}
	for i := 0; i < size; i++ {
		if i < size-1 {
			nodes[i].SetCollect(collect)
			nodes[i].SetNext(nodes[i+1])
		}
	}
	return &SingleCoroutineStage{
		s:         s,
		beginNode: nodes[0],
		ctx:       ctx,
	}, nil
}

func (mcs *SingleCoroutineStage) Run() {
	mcs.beginNode.Execute(mcs.s)
}

func (mcs *SingleCoroutineStage) getStat(statCh func(chan Stat)) {
	// statCh()1

}

type MultiCoroutineStage struct {
	p       *WrappedPool
	limiter Limiter
	ctx     context.Context
	errMsg  chan string
	se      *StatisticsEngine
}

func ConnectMultiCoroutineNode(t string, nodes []Node, num, per int, ctx context.Context, s *Shared) (*MultiCoroutineStage, error) {
	scs, err := ConnectSingleCoroutineNode(nodes, ctx, s, true)
	if err != nil {
		return nil, err
	}

	mcs := &MultiCoroutineStage{
		limiter: NewLimiter(t, per, ctx),
		ctx:     ctx,
		errMsg:  make(chan string, MessageThreshold),
	}
	se := NewStatisticalEngine()
	for _, n := range nodes {
		se.add(n.GetID(), ctx)
		n.SetStatCall(mcs.stat)
	}
	mcs.se = se
	p, err := NewPool(num, ctx, func() {
		scs.Run()
	})
	if err != nil {
		return nil, err
	}
	mcs.p = p
	mcs.run()
	return mcs, nil
}

func (mcs *MultiCoroutineStage) run() {
	mcs.p.RunByLimit(mcs.limiter)
}

func (mcs *MultiCoroutineStage) stat(id string, s Stat) {
	mcs.se.pushStat(id, s)
}
