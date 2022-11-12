package engine

import (
	"common"
	"context"
	"errors"
	"fmt"
	"plugin/shared"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

var ncPool sync.Pool

func acquireNodeContext(p Pipeline) *nodeContext {
	v := ncPool.Get()
	var nc *nodeContext
	if v == nil {
		nc = newNodeContext()
	} else {
		nc = v.(*nodeContext)
	}
	nc.initNodeContext(p)
	return nc
}

func releaseNodeContext(nc *nodeContext) {
	nc.reset()
	ncPool.Put(nc)
}

type nodeContext struct {
	statusList   map[string]chan uint8
	inputs       map[string][]Transport
	inputLen     map[string]int
	outputs      map[string]Transport
	taskId       string
	globalStatus chan uint8
	lock         chan uint8
}

func newNodeContext() *nodeContext {
	nc := &nodeContext{}
	nc.statusList = make(map[string]chan uint8)
	nc.inputs = make(map[string][]Transport)
	nc.inputLen = make(map[string]int)
	nc.outputs = make(map[string]Transport)
	nc.lock = make(chan uint8, 1)
	nc.lock <- 1
	nc.globalStatus = make(chan uint8, 1)
	return nc
}

func (nc *nodeContext) initNodeContext(p Pipeline) {
	for id, n := range p.Nodes {
		l := len(n.GetParents())
		nc.statusList[id] = make(chan uint8, 1)
		nc.inputs[id] = []Transport{}
		nc.inputLen[id] = l
		nc.taskId = p.TaskId
	}
}

func (nc *nodeContext) reset() {
	for _, v := range nc.statusList {
		close(v)
	}
	clearMap(nc.statusList)
	clearMap(nc.inputs)
	clearMap(nc.inputLen)
	clearMap(nc.outputs)
	nc.taskId = ""
	nc.lock = make(chan uint8, 1)
	nc.lock <- 1
	nc.globalStatus = make(chan uint8, 1)
}

func (nc *nodeContext) putOutput(nodeId string, t Transport) {
	<-nc.lock
	nc.outputs[nodeId] = t
	flag := true
	for k, _ := range nc.inputLen {
		if _, ok := nc.outputs[k]; !ok {
			flag = false
			break
		}
	}
	nc.lock <- 1
	if flag {
		nc.globalStatus <- 1
	}
}

func (nc *nodeContext) mergeTransport(nodeId string) (map[string]string, error) {
	<-nc.lock
	defer func() {
		nc.lock <- 1
	}()
	result := map[string]string{}
	keyMap := map[string]bool{}
	for _, t := range nc.inputs[nodeId] {
		if t.Status == common.Failed {
			return nil, errors.New(t.Err)
		}
		for k, v := range t.Output {
			if _, ok := keyMap[k]; ok {
				return nil, fmt.Errorf("nodId: '%s', receive the same key value", nodeId)
			}
			result[k] = v
		}
	}
	return result, nil
}

func (nc *nodeContext) start(nodeId string) {
	ch := nc.statusList[nodeId]
	if !common.IsClosed(ch) {
		<-nc.statusList[nodeId]
	}
}

func (nc *nodeContext) put(nodeId string, t Transport) {
	<-nc.lock
	defer func() {
		nc.lock <- 1
	}()
	nc.inputs[nodeId] = append(nc.inputs[nodeId], t)
	if len(nc.inputs[nodeId]) == nc.inputLen[nodeId] {
		nc.statusList[nodeId] <- 1
	}
}

type MapType interface {
	chan uint8 | []Transport | int | Transport
}

func clearMap[T MapType](data map[string]T) {
	for k := range data {
		delete(data, k)
	}
}

type Line struct {
	From string
	To   string
}

type Pipeline struct {
	Lines   []Line
	Nodes   map[string]Node
	Config  Config
	TaskId  string
	running int32
}

type Config struct {
	Type  string
	QPS   int
	Time  int
	Users int
	wp    *WrappedPool
	limit *Limiter
}

func NewPipeline(taskId string, c Config, ls []Line, ns map[string]Node) (*Pipeline, error) {
	p := &Pipeline{
		Lines:  ls,
		Nodes:  ns,
		Config: c,
		TaskId: taskId,
	}
	if err := p.check(); err != nil {
		return nil, err
	}
	p.connectNode()
	return p, nil
}

func (p *Pipeline) clear()  {
	for _, node := range p.Nodes {
		c := node.GetContainer()
		pluginId := c.pluginId
		version := c.version
		shared.PM.DeleteTask(pluginId, p.TaskId, version)
	}
}

func (p *Pipeline) check() error {
	if _, ok := common.PipelineTypeDec[p.Config.Type]; !ok {
		return errors.New("only support default, performance pipeline")
	}
	return nil
}

func (p *Pipeline) getContext(id string) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.Trace, id)
	if p.Config.Type == common.Performance {
		return context.WithTimeout(ctx, time.Duration(p.Config.Time)*time.Second)
	} else {
		return context.WithCancel(ctx)
	}
}

func (p *Pipeline) Run(ctx context.Context, receive func(t Transport)) error {
	if p.Config.Type == common.Performance {
		return p.runPerformance(ctx, receive)
	} else if p.Config.Type == common.Default {
		p.runDefault(receive)
	}
	return nil
}

func (p *Pipeline) ResetQPS(qps int) error {
	if p.Config.limit != nil {
		limit := rate.Limit(qps)
		if qps == 0 {
			limit = rate.Limit(-1)
		}
		p.Config.limit.l.SetLimit(limit)
		return nil
	}
	return errors.New("this mode cannot set qps")
}

func (p *Pipeline) GetUsers() (int, error) {
	if p.Config.limit != nil {
		return p.Config.wp.p.Cap(), nil
	}
	return 0, errors.New("this mode cannot get users")
}

func (p *Pipeline) ResetUsers(users int) error {
	if p.Config.wp != nil {
		p.Config.wp.p.Tune(users)
		return nil
	}
	return errors.New("this mode cannot set users")
}

func (p *Pipeline) GetRunningWorks() (int, error) {
	if p.Config.wp != nil {
		return int(atomic.LoadInt32(&p.running)), nil
	}
	return 0, errors.New("this mode cannot get works")
}

func (p *Pipeline) connectNode() {
	for _, line := range p.Lines {
		if node, ok := p.Nodes[line.To]; ok {
			node.PutParent(line.From)
		}
		if node, ok := p.Nodes[line.From]; ok {
			node.PutNext(line.To)
		}
	}
}

func (p *Pipeline) getConcurrentNode() (Node, error) {
	var ns []Node
	for _, n := range p.Nodes {
		if n.GetType() == common.NodeMulti {
			ns = append(ns, n)
		}
	}
	l := len(ns)
	if l == 0 {
		return nil, nil
	}
	if l > 1 {
		return nil, errors.New("there can only be one concurrent node")
	}
	return ns[0], nil
}

func (p *Pipeline) runPerformance(ctx context.Context, receive func(t Transport)) error {
	wp, err := NewPool(p.Config.Users, func() {
		p.run(receive)
	})
	p.Config.wp = wp
	if err != nil {
		return err
	}
	qps := p.Config.QPS
	if qps <= 0 {
		wp.Run(ctx)
	} else {
		l := NewLimiter(qps)
		p.Config.limit = &l
		wp.RunByLimit(ctx, l)
	}
	return nil
}

func (p *Pipeline) runDefault(receive func(t Transport)) {
	p.run(receive)
}

func (p *Pipeline) run(receive func(t Transport)) {
	atomic.AddInt32(&p.running, 1)
	defer atomic.AddInt32(&p.running, -1)
	nc := acquireNodeContext(*p)
	defer releaseNodeContext(nc)
	for _, n := range p.Nodes {
		n.Run(nc)
	}
	<-nc.globalStatus
	for _, t := range nc.outputs {
		receive(t)
	}
}
