package models

import (
	"common"
	"engine"
	"errors"
	"plugin/shared"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

type PipelineVO struct {
	Index      int    `json:"index,omitempty"`
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Describe   string `json:"describe,omitempty"`
	Type       string `json:"type,omitempty"`
	Config     Config `json:"config,omitempty"`
	CreateTime string `json:"createTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

func (pv *PipelineVO) GetEnginePipelineConfig() engine.Config {
	return engine.Config{
		Type:  pv.Type,
		QPS:   pv.Config.QPS,
		Time:  pv.Config.Time,
		Users: pv.Config.Users,
	}
}

type Config struct {
	QPS   int `json:"qps,omitempty"`
	Time  int `json:"time,omitempty"`
	Users int `json:"users,omitempty"`
}

type FlowVO struct {
	PipelineId string       `json:"pipelineId,omitempty"`
	Nodes      []FlowNodeVO `json:"nodes,omitempty"`
	Edges      []FlowEdgeVO `json:"edges,omitempty"`
}

func (fv *FlowVO) GetNodesLines() (map[string]engine.Node, []engine.Line, error) {
	var ls []engine.Line
	for _, e := range fv.Edges {
		ls = append(ls, e.getLine())
	}
	ns := map[string]engine.Node{}
	for _, n := range fv.Nodes {
		if n.Type == common.Plugin {
			pID := n.Data.PluginID
			version, err := shared.PM.GetNewestVersion(pID)
			if err != nil {
				return nil, nil, err
			}
			ns[n.ID] = engine.NewPluginNode(n.ID, n.Type, pID, version, n.Data.GetDefaultData())
		}
	}
	return ns, ls, nil
}

type FlowNodeVO struct {
	ID       string             `json:"id,omitempty"`
	NodeId   string             `json:"nodeId,omitempty"`
	Position map[string]float64 `json:"position,omitempty"`
	Type     string             `json:"type,omitempty"`
	Data     Data               `json:"data,omitempty"`
}

func (fnv *FlowNodeVO) getNode() (engine.Node, error) {
	var node engine.Node
	if fnv.Type == common.Plugin {
		pID := fnv.Data.PluginID
		version, err := shared.PM.GetNewestVersion(pID)
		if err != nil {
			return nil, err
		}
		node = engine.NewPluginNode(fnv.ID, fnv.Type, pID, version, fnv.Data.GetDefaultData())
	}
	return node, nil
}

type FlowEdgeVO struct {
	ID       string `json:"id,omitempty"`
	Source   string `json:"source,omitempty"`
	Target   string `json:"target,omitempty"`
	Type     string `json:"type,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

func (fev *FlowEdgeVO) getLine() engine.Line {
	return engine.Line{
		From: fev.Source,
		To:   fev.Target,
	}
}

type Data struct {
	DefaultInputs []Input `json:"defaultInputs,omitempty"`
	Describe      string  `json:"describe,omitempty"`
	Name          string  `json:"name,omitempty"`
	PluginID      string  `json:"pluginID,omitempty"`
	Code          string  `json:"code,omitempty"`
	Language      string  `json:"language,omitempty"`
}

func (d *Data) GetMultiNodeConfig() (map[string]int, error) {
	result := map[string]int{}
	for _, input := range d.DefaultInputs {
		if input.Name == common.MultiNodeQPS || input.Name == common.MultiNodeUsers || input.Name == common.MultiNodeTime {
			intVar, err := strconv.Atoi(input.Value)
			if err != nil {
				return nil, errors.New("concurrent node config value of type must be int")
			}
			result[input.Name] = intVar
		}
	}
	return result, nil
}

func (d *Data) GetDefaultData() map[string]string {
	r := map[string]string{}
	for _, item := range d.DefaultInputs {
		r[item.Name] = item.Value
	}
	return r
}

func (pv *PipelineVO) GetBaseBson() bson.M {
	base := bson.M{}
	if pv.Describe != "" {
		base["describe"] = pv.Describe
	}
	if pv.Type != "" {
		base["type"] = pv.Type
	}
	return base
}

type Pagination struct {
	Items interface{} `json:"items" bson:"items"`
	Total int64       `json:"total" bson:"total"`
}

type PipelineType struct {
	Name  string `json:"nane" bson:"items"`
	Value string `json:"value" bson:"items"`
}
