package core

type Node interface {
	Execute(interface{})
	SetNode(Node)
}

type HttpNode struct {
	next Node
}

func (hn *HttpNode) SetNext(node Node) {
	hn.next = node
}

func (hn *HttpNode) Execute(i interface{}) {
	hn.next.Execute(i)
}

type PostProcessorNode struct {
}

type ScriptNode struct {
}

type AssertNode struct {
}

type ReportNode struct {
}
