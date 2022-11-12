package db

import (
	"common"
	"errors"
	"time"
	"common/bindmodels"
	"web/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pipeline struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Describe   string             `bson:"describe,omitempty"`
	Type       string             `bson:"type,omitempty"`
	Status     string             `bson:"status,omitempty"`
	Config     Config             `bson:"config,omitempty"`
	CreateTime time.Time          `bson:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"updateTime,omitempty"`
	Flow       *Flow              `bson:"flow,omitempty"`
}

type Config struct {
	QPS   int `bson:"qps,omitempty"`
	Time  int `bson:"time,omitempty"`
	Users int `bson:"users,omitempty"`
}

type Flow struct {
	PipelineId string     `bson:"pipelineId,omitempty"`
	Nodes      []FlowNode `bson:"nodes,omitempty"`
	Edges      []FlowEdge `bson:"edges"`
}

type FlowNode struct {
	ID       string             `bson:"id,omitempty"`
	NodeId   string             `bson:"nodeId,omitempty"`
	Position map[string]float64 `bson:"position,omitempty"`
	Type     string             `bson:"type,omitempty"`
	Data     *Data              `bson:"data,omitempty"`
}

type FlowEdge struct {
	ID       string `bson:"id,omitempty"`
	Source   string `bson:"source,omitempty"`
	Target   string `bson:"target,omitempty"`
	Type     string `bson:"type,omitempty"`
	Animated bool   `bson:"animated,omitempty"`
}

type Data struct {
	DefaultInputs []Input `bson:"defaultInputs,omitempty"`
	Describe      string  `bson:"describe,omitempty"`
	Name          string  `bson:"name,omitempty"`
	PluginID      string  `bson:"pluginID,omitempty"`
	Code          string  `bson:"code,omitempty"`
	Type          string  `bson:"type,omitempty"`
	Language      string  `bson:"language,omitempty"`
}

func (data *Data) CompletePlugin(t string) error {
	if t == "" {
		return nil
	}
	if t != common.Plugin {
		return nil
	}
	if t == common.Plugin && data.PluginID == "" {
		return errors.New("请关联插件")
	}

	p, err := GetPluginByID(data.PluginID)
	if err != nil {
		return err
	}
	data.Code = p.Code
	data.Language = p.Language
	return nil
}

func UpdateFlow(flow *Flow) error {
	id, err := getIDFromInterface(flow.PipelineId)
	if err != nil {
		return err
	}
	p, err := GetPipelineByID(id)
	if err != nil {
		return err
	}
	for _, fn := range flow.Nodes {
		if err := fn.Data.CompletePlugin(fn.Type); err != nil {
			return err
		}
	}
	p.Flow = flow
	return UpdatePipeline(&p)
}

func GetFlowByPipelineId(pipelineId string) (*Flow, error) {
	id, err := getIDFromInterface(pipelineId)
	if err != nil {
		return nil, err
	}
	p, err := GetPipelineByID(id)
	if err != nil {
		return nil, err
	}
	return p.Flow, nil
}

func (p *Pipeline) GetPipelineVO() models.PipelineVO {
	return models.PipelineVO{
		ID:         p.ID.Hex(),
		Describe:   p.Describe,
		Name:       p.Name,
		Type:       p.Type,
		CreateTime: p.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: p.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}

func GetPipelines(q bindmodels.Query) ([]Pipeline, error) {
	filterOption, err := getFilterOption(q)
	filterOption["status"] = common.NormalStatus
	if err != nil {
		return []Pipeline{}, err
	}
	option := getPagingOption(q)
	return commonGet[Pipeline](common.PipelineCollection, filterOption, option)
}

func CreatePipeline(p *Pipeline) (primitive.ObjectID, error) {
	// b := pv.GetBaseBson()
	p.CreateTime = time.Now().In(common.GetLocalZone())
	p.UpdateTime = time.Now().In(common.GetLocalZone())
	p.Status = common.NormalStatus
	return commonCreate(p, common.PipelineCollection)
}

func UpdatePipeline(p *Pipeline) error {
	p.UpdateTime = time.Now().In(common.GetLocalZone())
	return commontUpdate(p, p.ID, common.PipelineCollection)
}

func GetPipelineCount() (int64, error) {
	ctx, cancel, c := getCollection(common.PipelineCollection)
	defer cancel()
	count, err := c.CountDocuments(ctx, getNormalStatusFilter())
	if err != nil {
		return 0, errors.New("get pipeline count error")
	}
	return count, nil
}

// GetPipelineByID 需要优化
func GetPipelineByID(id primitive.ObjectID) (Pipeline, error) {
	p, err := commontGetByID[Pipeline](id, getNormalStatusFilter(), common.PipelineCollection)
	if err != nil {
		return Pipeline{}, err
	}
	return p, nil
}

func DeletePipeline(id primitive.ObjectID) error {
	ctx, cancel, c := getCollection(common.PipelineCollection)
	defer cancel()
	update := bson.M{"$set": bson.M{"status": common.DeletedStatus}}
	if _, err := c.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
		return errors.New("delete pipeline failed")
	}
	return nil
}

func AddNode(id primitive.ObjectID, nodeId primitive.ObjectID) error {
	ctx, cancel, c := getCollection(common.PipelineCollection)
	defer cancel()
	if _, err := c.UpdateByID(ctx, id, bson.M{"$push": bson.M{"nodes": nodeId}}); err != nil {
		return errors.New("add node to pipeline error")
	}
	return nil
}

func DeleteNodeFromPipeline(id primitive.ObjectID, nodeId primitive.ObjectID) error {
	ctx, cancel, c := getCollection(common.PipelineCollection)
	defer cancel()
	if _, err := c.UpdateByID(ctx, id, bson.M{"$pull": bson.M{"nodes": nodeId}}); err != nil {
		return errors.New("delete node from pipeline error")
	}
	return nil
}

func getNormalStatusFilter() bson.M {
	return bson.M{"status": common.NormalStatus}
}
