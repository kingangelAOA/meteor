package core

import (
	"context"
	"errors"
)

type Stage interface {
	Run(p Shared)
}

func ConnectSingleCoroutineNode(nodes []Node, s Shared, ctx context.Context) (*SingleCoroutineStage, error) {
	size := len(nodes)
	if size == 0 {
		return nil, errors.New("must be at least one node")
	}
	for i := 0; i < size; i++ {
		if i < size-1 {
			nodes[i].SetNode(nodes[i+1])
		}
	}
	return &SingleCoroutineStage{
		s:         s,
		beginNode: nodes[0],
		ctx:       ctx,
	}, nil
}

func ConnectMultiCoroutineNode(t string, nodes []Node, num, per int, ctx context.Context, s Shared) (*MultiCoroutineStage, error) {
	scs, err := ConnectSingleCoroutineNode(nodes, s, ctx)
	if err != nil {
		return nil, err
	}
	mcs := &MultiCoroutineStage{
		limiter: NewLimiter(t, per, ctx),
		ctx:     ctx,
		errMsg:  make(chan string, MessageThreshold),
	}
	p, err := NewPool(num, ctx, func() {
		scs.Run()
	})
	if err != nil {
		return nil, err
	}
	mcs.p = p
	return mcs, nil
}

type SingleCoroutineStage struct {
	s         Shared
	beginNode Node
	ctx       context.Context
}

func (mcs *SingleCoroutineStage) Run() error {
	return mcs.beginNode.Execute(&mcs.s)
}

type MultiCoroutineStage struct {
	p       *WrappedPool
	limiter Limiter
	ctx     context.Context
	errMsg  chan string
}

func (scs *MultiCoroutineStage) Run(s Shared) {
	scs.p.RunByLimit(scs.limiter, s)
}
