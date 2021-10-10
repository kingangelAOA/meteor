package core

const (
	NodeErrMsgThreshold = 1000
)

type Node interface {
	Execute(*Shared)
	SetNext(Node)
	SetStatCall(func(string, Stat))
	SetCollect(bool)
	GetID() string
	GetErrMsg() string
}

type BaseNode struct {
	next    Node
	collect bool
	id      string
	call    func(string, Stat)
	errMsg  chan string
}

func NewBaseNode(id string) BaseNode {
	return BaseNode{
		id:     id,
		errMsg: make(chan string, NodeErrMsgThreshold),
	}
}

func (bn *BaseNode) putErrMsg(e string) {
	if len(bn.errMsg) > 1000 {
		<-bn.errMsg
	}
	bn.errMsg <- e
}

func (bn *BaseNode) SetStatCall(f func(string, Stat)) {
	bn.call = f
}

func (bn *BaseNode) GetID() string {
	return bn.id
}

func (bn *BaseNode) GetErrMsg() string {
	if len(bn.errMsg) > 0 {
		return <-bn.errMsg
	}
	return ""
}

func (bn *BaseNode) SetCollect(s bool) {
	bn.collect = s
}

func (bn *BaseNode) SetNext(node Node) {
	bn.next = node
}

type HttpNode struct {
	WrappedRequest WrappedRequest
	http           *Http
	BaseNode
}

func (hn *HttpNode) Execute(s *Shared) error {
	// r, err := hn.WrappedRequest.GetRequest(s)
	// if err != nil {
	//  return err
	// }
	// NewHttpMessage()
	// hn.http.Execute()
	// hn.next.Execute()
	return nil
}

type ScriptNode struct {
	Name string
	Type string
	BaseNode
}

func NewScript(id, name, t string) *ScriptNode {
	return &ScriptNode{
		Name:     name,
		Type:     t,
		BaseNode: NewBaseNode(id),
	}
}

func (sn *ScriptNode) Execute(s *Shared) {
	k, err := GetScriptKey(sn.id, sn.Name)
	if err != nil {
		sn.putErrMsg(err.Error())
	}
	tm := AcquireScriptMessage(k, TengoType, CopyMap(s.Data))
	DefaultScriptService.PutMessage(tm)
	ok := <-tm.Ok
	if !ok {
		sn.putErrMsg(tm.ErrMsg)
	}
	s.SetData(tm.Data)
	stat := tm.GetStat()
	stat.ok = ok
	if sn.collect {
		sn.call(sn.id, stat)
	}
	ReleaseScriptMessage(tm)
	if sn.next != nil && ok {
		sn.next.Execute(s)
	}
}

type PostProcessorNode struct {
}

type AssertNode struct {
}

type ReportNode struct {
}
