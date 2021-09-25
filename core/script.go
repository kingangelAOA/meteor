package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/panjf2000/ants/v2"
)

const (
	ScriptCtx       = "ctx"
	ScriptPrints    = "prints"
	ErrMsgThreshold = 1000
	ExpiryDuration  = 60
)

var (
	ErrorNewScriptServerPool = errors.New("new script pool error")
)

var (
	TengoVMPool      map[string]*sync.Pool
	TengoMessagePool *sync.Pool
)

func init() {
	TengoVMPool = make(map[string]*sync.Pool)
	TengoMessagePool = &sync.Pool{}
}

type ScriptService struct {
	BaseService
	script Script
}

func NewScriptService(s Script, num int, ctx context.Context) (*ScriptService, error) {
	ss := &ScriptService{
		script: s,
	}
	p, err := ants.NewPool(num, ants.WithPanicHandler(func(p interface{}) {
		ss.PutErrMsg(fmt.Sprintf("worker exits from a panic: %v\n", p))
	}), ants.WithExpiryDuration(ExpiryDuration*time.Minute))
	if err != nil {
		return nil, ErrorNewScriptServerPool
	}
	bs := NewBaseService(p, ctx)
	ss.BaseService = *bs
	return ss, nil
}

func (ss *ScriptService) Run() {
	go ss.messageQueue()
}

func (ss *ScriptService) SetGoNum(size int) {
	ss.p.Tune(size)
}

func (ss *ScriptService) PutErrMsg(m string) {
	if len(ss.ErrMsg) > ErrMsgThreshold {
		<-ss.ErrMsg
	}
	ss.ErrMsg <- m
}

func (ss *ScriptService) messageQueue() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			ss.Code = 500
			ss.PutErrMsg(err.Error())
		}
	}()
	for {
		select {
		case <-ss.ctx.Done():
			return
		case pm := <-ss.ms:
			err := ss.p.Submit(func() {
				ss.script.Execute(pm)
			})
			if err != nil {
				pm.SetErr(err.Error())
			}
		}
	}
}

type TengoMessage struct {
	BaseMessage
	Name string
}

func NewTengoMessage(name string, s *Shared) *TengoMessage {
	return &TengoMessage{
		Name: name,
		BaseMessage: BaseMessage{
			s:  s,
			Ok: make(chan bool, 1),
		},
	}
}

func (tm *TengoMessage) Reset() {
	tm.Name = ""
	if len(tm.Ok) == 1 {
		<-tm.Ok
	}
}

func (tm *TengoMessage) GetName() string {
	return tm.Name
}

func AcquireTengoMessage(name string, s *Shared) *TengoMessage {
	v := TengoMessagePool.Get()
	if v == nil {
		return NewTengoMessage(name, s)
	}
	tm := v.(*TengoMessage)
	tm.Name = name
	tm.s = s
	return tm
}

func ReleaseTengoMessage(tm *TengoMessage) {
	tm.Reset()
	TengoMessagePool.Put(tm)
}

type Script interface {
	Execute(Message)
	AddScript(string, string)
}

type TengoScript struct {
	scriptMap map[string]string
	ctx       context.Context
}

type TengoVM struct {
	name   string
	script *tengo.Script
}

func (tvm *TengoVM) Reset() {
}

func (tvm *TengoVM) Run(m Message, ctx context.Context) {
	tvm.script.Add(ScriptCtx, m.GetShared().Data)
	tvm.script.Add(ScriptPrints, []interface{}{})
	compiled, err := tvm.script.RunContext(ctx)
	errMsg := []string{}
	if err != nil {
		errMsg = append(errMsg, fmt.Sprintf("tengo script run error: %s", err.Error()))
	}
	if nil != compiled {
		if v, ok := compiled.Get(ScriptCtx).Value().(map[string]interface{}); ok {
			m.SetShared(v)
		} else {
			errMsg = append(errMsg, "tengo script run error: ctx is not map")
		}

		if out, ok := compiled.Get(ScriptPrints).Value().([]interface{}); ok {
			buffer := bytes.Buffer{}
			for _, v := range out {
				buffer.WriteString(v.(string))
				buffer.WriteString("\n")
			}
			m.SetPrints(buffer.String())
		} else {
			errMsg = append(errMsg, "tengo script run error: prints is not slice")
		}
	}

	if len(errMsg) > 0 {
		buffer := bytes.Buffer{}
		for _, v := range errMsg {
			buffer.WriteString(v)
			buffer.WriteString("\n")
		}
		m.SetErr(buffer.String())
	} else {
		m.SetOk(true)
	}
}

func AcquireTengoVM(name, code string) *TengoVM {
	if pool, ok := TengoVMPool[name]; ok {
		v := pool.Get()
		if v == nil {
			return &TengoVM{
				name:   name,
				script: newTengoVM(code),
			}
		}
		return v.(*TengoVM)
	} else {
		TengoVMPool[name] = &sync.Pool{}
		return &TengoVM{
			name:   name,
			script: newTengoVM(code),
		}
	}
}

func ReleaseTengoVM(tvm *TengoVM) {
	tvm.Reset()
	if v, ok := TengoVMPool[tvm.name]; ok {
		v.Put(tvm)
	} else {
		pool := &sync.Pool{}
		pool.Put(tvm)
		TengoVMPool[tvm.name] = pool
	}
}

func newTengoVM(code string) *tengo.Script {
	script := tengo.NewScript([]byte(code))
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	return script
}

func NewTengoScript(ctx context.Context) *TengoScript {
	return &TengoScript{
		scriptMap: map[string]string{},
		ctx:       ctx,
	}
}

func (ts *TengoScript) AddScript(name, code string) {
	ts.scriptMap[name] = code
}

func (ts *TengoScript) Execute(m Message) {
	if code, ok := ts.scriptMap[m.GetName()]; ok {
		tvm := AcquireTengoVM(m.GetName(), code)
		tvm.Run(m, ts.ctx)
		ReleaseTengoVM(tvm)
	} else {
		m.SetErr(fmt.Sprintf("script '%s' is not exist", m.GetName()))
	}
}