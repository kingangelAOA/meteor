package core

import "context"

type Stage interface {
	Run(p *Shared) *Shared
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

type SingleCoroutineStage struct {
	beginNode Node
}

func (mcs *SingleCoroutineStage) Run(s *Shared) {
	mcs.beginNode.Execute(s)
}

type MultiCoroutineStage struct {
	pool    WrappedPool
	limiter Limiter
}

func (scs *MultiCoroutineStage) Run(s *Shared) {
	scs.pool.RunByLimit(scs.limiter, s)
}
