package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/panjf2000/ants/v2"
)

const (
	ScriptCtx        = "ctx"
	ScriptPrints     = "prints"
	MessageThreshold = 5000
)

var (
	ErrorNewScriptServerPool = errors.New("new script pool error")
)

var (
	TengoVMPool      map[string]sync.Pool
	TengoMessagePool sync.Pool
)

func init() {
	TengoVMPool = make(map[string]sync.Pool)
}

type ScriptService struct {
	script   Script
	ctx      context.Context
	Code     int
	p        *ants.Pool
	ErrMsg   string
	messages chan Message
}

func NewScriptService(s Script, num int, ctx context.Context) (*ScriptService, error) {
	p, err := ants.NewPool(num, ants.WithPreAlloc(true))
	if err != nil {
		return nil, ErrorNewScriptServerPool
	}
	return &ScriptService{
		script:   s,
		ctx:      ctx,
		p:        p,
		messages: make(chan Message, MessageThreshold),
	}, nil
}

func (ss *ScriptService) Run() {
	go ss.messageQueue(ss.messages)
}

func (ss *ScriptService) messageQueue(mCh chan Message) {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			ss.Code = 500
			ss.ErrMsg = err.Error()
		}
	}()
	for {
		select {
		case <-ss.ctx.Done():
			return
		case pm := <-mCh:
			err := ss.p.Submit(func() {
				ss.script.Execute(pm)
			})
			if err != nil {
				pm.SetErr(err.Error())
			}
		}
	}
}

func (ss *ScriptService) PutMessage(m Message) {
	ss.messages <- m
}

type Message interface {
	GetName() string
	GetWr() *WrappedContext
	SetWr(map[string]interface{})
	SetErr(string)
	SetOk(bool)
	GetPrings() string
	SetPrints(string)
}

type TengoMessage struct {
	Name   string
	WrCh   *WrappedContext
	Ok     chan bool
	ErrMsg string
	Prints string
}

func AcquireTengoMessage(name string, WrCh *WrappedContext) *TengoMessage {
	v := TengoMessagePool.Get()
	if v == nil {
		return &TengoMessage{
			Name: name,
			WrCh: WrCh,
			Ok:   make(chan bool, 1),
		}
	}
	tm := v.(*TengoMessage)
	tm.Name = name
	tm.WrCh = WrCh
	return tm
}

func ReleaseTengoMessage(tm *TengoMessage) {
	tm.Reset()
	TengoMessagePool.Put(tm)
}

func (tm *TengoMessage) Reset() {
	tm.Name = ""
	tm.WrCh.Reset()
	if len(tm.Ok) == 1 {
		<-tm.Ok
	}
}

func (tm *TengoMessage) GetName() string {
	return tm.Name
}

func (tm *TengoMessage) GetWr() *WrappedContext {
	return tm.WrCh
}

func (tm *TengoMessage) SetWr(data map[string]interface{}) {
	tm.WrCh.Data = data
}

func (tm *TengoMessage) SetErr(msg string) {
	tm.ErrMsg = msg
}

func (tm *TengoMessage) SetOk(ok bool) {
	tm.Ok <- ok
}

func (tm *TengoMessage) GetPrings() string {
	return tm.Prints
}

func (tm *TengoMessage) SetPrints(p string) {
	tm.Prints = p
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
	tvm.script.Add(ScriptCtx, m.GetWr().Data)
	tvm.script.Add(ScriptPrints, []interface{}{})
	compiled, err := tvm.script.RunContext(ctx)
	errMsg := []string{}
	if err != nil {
		errMsg = append(errMsg, fmt.Sprintf("tengo script run error: %s", err.Error()))
	}
	if nil != compiled {
		if v, ok := compiled.Get(ScriptCtx).Value().(map[string]interface{}); ok {
			m.SetWr(v)
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
		m.SetOk(false)
	} else {
		m.SetOk(true)
	}
	ReleaseTengoVM(tvm)
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
		TengoVMPool[name] = sync.Pool{}
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
		pool := sync.Pool{}
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
	} else {
		m.SetErr(fmt.Sprintf("script '%s' is not exist", m.GetName()))
		m.SetOk(false)
	}
}
