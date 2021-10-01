package core

type Node interface {
	Execute(*Shared) error
	SetNode(Node)
}

type HttpNode struct {
	WrappedRequest WrappedRequest
	http           *Http
	next           Node
}

func (hn *HttpNode) SetNext(node Node) {
	hn.next = node
}

func (hn *HttpNode) Execute(s *Shared) error {
	// r, err := hn.WrappedRequest.GetRequest(s)
	// if err != nil {
	// 	return err
	// }
	// NewHttpMessage()
	// hn.http.Execute()
	// hn.next.Execute()
	return nil
}

type PostProcessorNode struct {
}

type ScriptNode struct {
}

type AssertNode struct {
}

type ReportNode struct {
}
