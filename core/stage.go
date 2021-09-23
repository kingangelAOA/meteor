package core

import "context"

type Stage interface {
	Run(p *WrappedContext) *WrappedContext
}

func ConnectSingleCoroutineNode(nodes []Node) Node {
	size := len(nodes)
	for i := 0; i < size; i++ {
		if i < size-1 {
			nodes[i].SetNode(nodes[i+1])
		}
	}
	return nodes[0]
}

func ConnectMultiCoroutineNode(nodes []Node, num int, ctx context.Context) (*WrappedPool, error) {
	node := ConnectSingleCoroutineNode(nodes)
	return NewPool(num, ctx, node.Execute)
}

type MultiCoroutineStage struct {
	beginNode Node
}

func (mcs *MultiCoroutineStage) Run(wc *WrappedContext) {
	mcs.beginNode.Execute(wc)
}

type SingleCoroutineStage struct {
	pool    WrappedPool
	limiter Limiter
}

func (scs *SingleCoroutineStage) Run(wc *WrappedContext) {
	scs.pool.RunByLimit(scs.limiter, wc)
}
