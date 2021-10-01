package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

// func init() {
// 	TengoVMPool = make(map[string]*sync.Pool)
// 	TengoMessagePool = &sync.Pool{}
// }

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
		Name:        name,
		BaseMessage: NewBaseMessage(s),
	}
}

func (tm *TengoMessage) Reset() *Shared {
	tm.Name = ""
	return tm.BaseMessage.reset()
}

func (tm *TengoMessage) GetName() string {
	return tm.Name
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
	// f := os.NewFile(1, "cache")
	// old := os.Stdout
	// os.Stdout = f
	compiled, err := tvm.script.RunContext(ctx)
	// f.Sync()
	// os.Stdout = old
	// b, _ := ioutil.ReadAll(f)
	// fmt.Println("*******", f)
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
				s, err := ToString(v)
				if err != nil {
					errMsg = append(errMsg, fmt.Sprintf("tengo script run error: %s", err.Error()))
				}
				buffer.WriteString(s)
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
