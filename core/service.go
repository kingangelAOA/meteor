package core

import (
	"context"
	"time"

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
	GetKey() string
	SetData(map[string]interface{})
	GetData() map[string]interface{}
	Reset()
	TimeCost() func()
	GetStat() Stat
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
	Data    map[string]interface{}
	Ok      chan bool
	endTime time.Time
	elapsed int
	ErrMsg  string
	Prints  string
}

func NewBaseMessage(data map[string]interface{}) BaseMessage {
	return BaseMessage{
		Data: data,
		Ok:   make(chan bool, 1),
	}
}

func (bm *BaseMessage) GetStat() Stat {
	return Stat{
		endTime: bm.endTime,
		elapsed: bm.elapsed,
	}
}

func (bm *BaseMessage) TimeCost() func() {
	begin := time.Now()
	return func() {
		end := time.Now()
		bm.elapsed = int(end.UnixMilli() - begin.UnixMilli())
		bm.endTime = end
	}
}

func (bm *BaseMessage) reset() {
	bm.ErrMsg = ""
	bm.Prints = ""
	for {
		if len(bm.Ok) > 0 {
			<-bm.Ok
		} else {
			break
		}
	}
	ClearMap(bm.Data)
}

func (bm *BaseMessage) SetData(data map[string]interface{}) {
	for k, v := range data {
		bm.Data[k] = v
	}
}

func (bm *BaseMessage) GetData() map[string]interface{} {
	return bm.Data
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
