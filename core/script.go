package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"meteor/models"
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
	TengoType       = "TENGO"
)

var (
	ErrorNewScriptServerPool = errors.New("new script pool error")
	DefaultScriptService     *ScriptService
)

type ScriptService struct {
	BaseService
	scripts map[string]Script
}

func NewScriptService(nodes models.ScriptNode, num int, ctx context.Context) (*ScriptService, error) {
	scripts, err := getMapScripts(nodes, ctx)
	if err != nil {
		return nil, err
	}
	ss := &ScriptService{
		scripts: scripts,
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

func (ss *ScriptService) getScript(name, t string) (Script, error) {
	script, ok := ss.scripts[t]
	if !ok {
		return nil, fmt.Errorf("no script of this type '%s' was found", t)
	}
	return script, nil
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
			sm := pm.(*ScriptMessage)
			s, err := ss.getScript(sm.Key, sm.Type)
			if err != nil {
				pm.SetErr(err.Error())
			} else {
				err := ss.p.Submit(func() {
					s.Execute(pm)
				})
				if err != nil {
					pm.SetErr(err.Error())
				}
			}
		}
	}
}

type ScriptMessage struct {
	BaseMessage
	Key  string
	Type string
}

func NewScriptMessage(key, t string, data map[string]interface{}) *ScriptMessage {
	return &ScriptMessage{
		Key:         key,
		Type:        t,
		BaseMessage: NewBaseMessage(data),
	}
}

func (tm *ScriptMessage) Reset() {
	tm.Key = ""
	tm.Type = ""
	tm.BaseMessage.reset()
}

func (tm *ScriptMessage) GetKey() string {
	return tm.Key
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
	key    string
	script *tengo.Script
}

func (tvm *TengoVM) Reset() {
}

func (tvm *TengoVM) Run(m Message, ctx context.Context) {
	defer m.TimeCost()
	tvm.script.Add(ScriptCtx, m.GetData())
	tvm.script.Add(ScriptPrints, []interface{}{})
	compiled, err := tvm.script.RunContext(ctx)
	errMsg := []string{}
	if err != nil {
		errMsg = append(errMsg, fmt.Sprintf("tengo script run error: %s", err.Error()))
	}
	if nil != compiled {
		if v, ok := compiled.Get(ScriptCtx).Value().(map[string]interface{}); ok {
			m.SetData(v)
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
	key := m.GetKey()
	if code, ok := ts.scriptMap[key]; ok {
		tvm := AcquireTengoVM(key, code)
		tvm.Run(m, ts.ctx)
		ReleaseTengoVM(tvm)
	} else {
		m.SetErr(fmt.Sprintf("script '%s' is not exist", key))
	}
}

func getMapScripts(nodes models.ScriptNode, ctx context.Context) (map[string]Script, error) {
	ms := make(map[string]Script)
	for k, bNodes := range nodes {
		if k == TengoType {
			nts := NewTengoScript(ctx)
			for _, node := range bNodes {
				k, err := GetScriptKey(node.ID, node.Name)
				if err != nil {
					return nil, err
				}
				nts.AddScript(k, node.Code)
			}
			ms[k] = nts
		}
	}
	return ms, nil
}
