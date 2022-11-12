package db

import (
	"common"
	"common/bindmodels"
	"errors"
	"mongoutil"
	"time"
	"web/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Describe   string             `bson:"describe,omitempty"`
	Type       string             `bson:"type,omitempty"`
	PluginID   primitive.ObjectID `bson:"pluginID,omitempty"`
	Plugin     *Plugin            `bson:"plugin,omitempty"`
	CreateTime time.Time          `bson:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"updateTime,omitempty"`
}

type FormConfig struct {
	Type          string
	Key           string
	Label         string
	SelectOptions []SelectOption
}

type SelectOption struct {
	Label string
	Value string
}

func (n *Node) GetNodeVO() models.NodeVO {
	var inputs []models.Input
	if nil != n.Plugin {
		for _, input := range n.Plugin.Inputs {
			inputs = append(inputs, models.Input{
				Name:  input.Name,
				Desc:  input.Desc,
				Value: input.Value,
			})
		}
	}

	nv := models.NodeVO{
		ID:            n.ID.Hex(),
		Describe:      n.Describe,
		Name:          n.Name,
		Type:          n.Type,
		DefaultInputs: inputs,
		PluginID:      n.PluginID.Hex(),
		CreateTime:    n.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTIme:    n.UpdateTime.Format("2006-01-02 15:04:05"),
	}
	completeSpecialNode(&nv)
	return nv
}

func GetNodes(q bindmodels.Query) ([]Node, error) {
	filterOption, err := getFilterOption(q)
	if err != nil {
		return []Node{}, err
	}
	option := getPagingOption(q)
	return commonGet[Node](common.NodeCollection, filterOption, option)
}

func GetAllNodes() ([]Node, error) {
	findOptions := &options.FindOptions{}
	return commonGet[Node](common.NodeCollection, bson.M{}, findOptions)
}

func GetBindingPluginNodes(q bindmodels.Query) ([]Node, error) {
	kv, err := parseSearch(q.Search)
	if err != nil {
		return []Node{}, err
	}
	lookup := &mongoutil.Lookup{
		From:         "plugin",
		LocalField:   "pluginID",
		ForeignField: "_id",
		As:           "plugin",
	}
	match := &mongoutil.Match{
		Querys: []mongoutil.Query{
			&mongoutil.EqualityCondition{
				Key: kv[0],
				Value: primitive.Regex{Pattern: kv[1], Options: ""},
			},
		},
	}
	unwind := &mongoutil.Unwind{
		Path: "$plugin",
		PreserveNullAndEmptyArrays: true,
	}
	dns, err := loadDefaultNodes()
	if err != nil {
		return nil, err
	}
	ns, err := commonAggregate[Node](common.NodeCollection, []bson.M{match.GetBson(), lookup.GetBson(), unwind.GetBson()}, getAggregateOpt())
	ns = append(ns, dns...)
	return ns, err
}

func GetNode(id primitive.ObjectID) (*Node, error) {
	ctx, cancel, c := getCollection(common.NodeCollection)
	defer cancel()
	sr := c.FindOne(ctx, bson.M{"_id": id})

	var n Node
	if err := sr.Decode(&n); err != nil {
		return nil, errors.New("decode nodes error")
	}
	return &n, nil
}

func CreateNode(n *Node) (primitive.ObjectID, error) {
	n.CreateTime = time.Now().In(common.GetLocalZone())
	n.UpdateTime = time.Now().In(common.GetLocalZone())
	return commonCreate(n, common.NodeCollection)
}

func UpdateNode(n *Node) error {
	n.UpdateTime = time.Now().In(common.GetLocalZone())
	return commontUpdate(n, n.ID, common.NodeCollection)
}

func GetNodeCount() (int64, error) {
	ctx, cancel, c := getCollection(common.NodeCollection)
	defer cancel()
	count, err := c.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, errors.New("get node count error")
	}
	return count, nil
}

func UpdatePluginFromNode(id primitive.ObjectID) error {
	ctx, cancel, c := getCollection(common.NodeCollection)
	defer cancel()
	if _, err := c.UpdateByID(ctx, id, bson.M{"pluginID": id}); err != nil {
		return errors.New("add node to pipeline error")
	}
	return nil
}

func DeleteNode(id string) error {
	ctx, cancel, c := getCollection(common.NodeCollection)
	defer cancel()
	objectID, err := ConverObjectID(id)
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.New("delete node error")
	}
	return nil
}

func loadDefaultNodes() ([]Node, error) {
	ido, err := getIDFromInterface(common.MGoroutingID)
	if err != nil {
		return nil, err
	}
	mgorouting := Node{
		ID:       ido,
		Name:     "并发节点",
		Describe: "作为整个流程的父节点存在",
		Type:     common.NodeMulti,
	}
	ns := []Node{}
	ns = append(ns, mgorouting)
	return ns, nil
}

func completeSpecialNode(ns *models.NodeVO) {
	if ns.Type == common.NodeMulti {
		ns.DefaultInputs = []models.Input{
			{
				Name:     common.MultiNodeQPS,
				Desc:     "设置稳定qps",
				Value:    "1",
				Required: true,
			},
			{
				Name:     common.MultiNodeUsers,
				Desc:     "用户数量",
				Value:    "1",
				Required: true,
			},
			{
				Name:     common.MultiNodeTime,
				Desc:     "运行时间(s)",
				Value:    "1",
				Required: true,
			},
		}
	}
}
