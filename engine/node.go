package engine

import (
	"common"
	"fmt"
	"plugin/proto"
	"plugin/shared"
	"time"
)

type pluginOutput struct {
	output *proto.OutPut
	error  error
}

type Node interface {
	Run(nc *nodeContext)
	PutParent(p string)
	PutNext(n string)
	GetParents() []string
	GetNexts() []string
	GetType() string
	GetContainer() Container
}

type baseNode struct {
	ID      string
	Type    string
	parents []string
	nexts   []string
}

func newBaseNode(id, t string) *baseNode {
	return &baseNode{
		ID:   id,
		Type: t,
	}
}

func (bn *baseNode) GetParents() []string {
	return bn.parents
}

func (bn *baseNode) GetType() string {
	return bn.Type
}

func (bn *baseNode) PutParent(p string) {
	bn.parents = append(bn.parents, p)
}

func (bn *baseNode) GetNexts() []string {
	return bn.nexts
}

func (bn *baseNode) PutNext(n string) {
	bn.nexts = append(bn.nexts, n)
}

type PluginNode struct {
	baseNode
	Container Container
}

func NewPluginNode(id, t, pluginId, version string, di map[string]string) *PluginNode {
	return &PluginNode{
		baseNode:  *newBaseNode(id, t),
		Container: NewContainer(pluginId, version, di),
	}
}

func (pn *PluginNode) GetContainer() Container {
	return pn.Container
}

func (pn *PluginNode) Run(nc *nodeContext) {
	go func(nc *nodeContext, pn *PluginNode) {
		id := pn.ID
		if len(pn.parents) > 0 {
			nc.start(id)
		}
		var output Transport
		result, err := nc.mergeTransport(id)
		if err != nil {
			output = NewErrTransport("parent node error", 0, time.Now().UnixMilli())
		} else {
			output = pn.Container.Run(result, nc.taskId)
		}
		for _, n := range pn.nexts {
			nc.put(n, output)
		}
		output.ID = id
		nc.putOutput(pn.ID, output)
	}(nc, pn)
}

type Container struct {
	pluginId     string
	version      string
	defaultInput map[string]string
}

func NewContainer(pluginId, version string, d map[string]string) Container {
	return Container{
		pluginId:     pluginId,
		version:      version,
		defaultInput: d,
	}
}

func (dc Container) isDefaultKey(key string) bool {
	if _, ok := dc.defaultInput[key]; ok {
		return true
	}
	return false
}

func (dc Container) Run(inputs map[string]string, taskId string) Transport {
	result := map[string]string{}
	for k, v := range dc.defaultInput {
		result[k] = v
	}
	for k, v := range inputs {
		result[k] = v
	}
	plugin, err := shared.PM.AcquirePlugin(dc.pluginId, dc.version, taskId)
	defer func() {
		if err := shared.PM.ReleasePlugin(dc.pluginId, dc.version, plugin); err != nil {
			fmt.Println(err.Error())
		}
	}()
	begin := time.Now()
	if err != nil {
		return NewErrTransport(err.Error(), 0, begin.UnixMilli())
	}
	poCh := make(chan pluginOutput, 1)
	go func() {
		po := pluginOutput{}
		output, err := plugin.Run(result)
		po.output = output
		po.error = err
		poCh <- po
	}()
	select {
	case po := <-poCh:
		tc := time.Since(begin)
		if po.error != nil {
			return NewErrTransport(po.error.Error(), int(tc.Milliseconds()), begin.UnixMilli())
		}
		if po.output.Rt <= 0 {
			po.output.Rt = int32(tc.Milliseconds())
		}
		if po.output.StartTime <= 0 {
			po.output.StartTime = begin.UnixMilli()
		}
		return NewDefaultTransport(dc.defaultInput, po.output.Data, int(po.output.Rt), po.output.StartTime)
	case <-time.After(common.PluginTimeout * time.Second):
		tc := time.Since(begin)
		errMsg := fmt.Sprintf("plugin timeout, timeout is '%d' seconds", common.PluginTimeout)
		return NewErrTransport(errMsg, int(tc.Seconds()), begin.UnixMilli())
	}
}
