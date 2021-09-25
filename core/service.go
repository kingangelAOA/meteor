package core

import (
	"context"

	"github.com/panjf2000/ants/v2"
)

const (
	MessageThreshold = 5000
)

type Service interface {
	Run()
	PutMessage(Message)
	SetGoNum(int)
	PutErrMsg(string)
}

type Message interface {
	GetName() string
	GetShared() *Shared
	SetShared(map[string]interface{})
	SetErr(string)
	SetOk(bool)
	GetPrints() string
	SetPrints(string)
}

type BaseService struct {
	p      *ants.Pool
	ctx    context.Context
	ErrMsg chan string
	Code   int
	ms     chan Message
}

func NewBaseService(p *ants.Pool, ctx context.Context) *BaseService {
	return &BaseService{
		p:      p,
		ctx:    ctx,
		ErrMsg: make(chan string, ErrMsgThreshold),
		ms:     make(chan Message, MessageThreshold),
	}
}

func (bs *BaseService) SetGoNum(size int) {
	bs.p.Tune(size)
}

func (bs *BaseService) PutErrMsg(m string) {
	if len(bs.ErrMsg) > ErrMsgThreshold {
		<-bs.ErrMsg
	}
	bs.ErrMsg <- m
}

func (ss *BaseService) PutMessage(m Message) {
	ss.ms <- m
}

type BaseMessage struct {
	s      *Shared
	Ok     chan bool
	ErrMsg string
	Prints string
}

func (bm *BaseMessage) GetShared() *Shared {
	return bm.s
}

func (bm *BaseMessage) SetShared(data map[string]interface{}) {
	for k, v := range data {
		bm.s.Set(k, v)
	}
}

func (bm *BaseMessage) SetErr(msg string) {
	bm.ErrMsg = msg
	bm.Ok <- false
}

func (bm *BaseMessage) SetOk(ok bool) {
	bm.Ok <- ok
}

func (tm *BaseMessage) GetPrints() string {
	return tm.Prints
}

func (tm *BaseMessage) SetPrints(p string) {
	tm.Prints = p
}
