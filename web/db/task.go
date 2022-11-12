package db

import (
	"common"
	"common/bindmodels"
	"encoding/json"
	"errors"
	"fmt"
	"mongoutil"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var tps = []string{common.PointTypeRT99, common.PointTypeRT95, common.PointTypeRT90, common.PointTypeRTAverage, common.PointTypeErrRate, common.PointTypeQPS}

type Task struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	PipelineId primitive.ObjectID   `bson:"pipelineId,omitempty"`
	Pipeline   Pipeline             `bson:"pipeline,omitempty"`
	Type       string               `bson:"type,omitempty"`
	Resutls    []NodeResult         `bson:"results,omitempty"`
	Points     []primitive.ObjectID `bson:"points,omitempty"`
	CreateTime time.Time            `bson:"createTime,omitempty"`
	UpdateTime time.Time            `bson:"updateTime,omitempty"`
}

func (t *Task) GetStatus() string {
	for _, r := range t.Resutls {
		if r.Status == common.Failed {
			return common.Failed
		}
	}
	return common.Success
}

type Point struct {
	Timestamp primitive.DateTime `bson:"timestamp,omitempty"`
	Metadata  Metadata           `bson:"metadata,omitempty"`
	Temp      float64            `bson:"temp,omitempty"`
}

type Metadata struct {
	ID   string `bson:"id,omitempty"`
	Type string `bson:"type,omitempty"`
}

type TypePoint map[string][]Point

type MonitorData struct {
	ID   string
	Name string
	TP   TypePoint
}

func (t *Task) GetConsume() int {
	consume := 0
	for _, r := range t.Resutls {
		consume += r.RT
	}
	return consume
}

type NodeResult struct {
	ID     string            `bson:"id,omitempty"`
	Name   string            `bson:"name,omitempty"`
	Err    string            `bson:"err,omitempty"`
	Status string            `bson:"status,omitempty"`
	RT     int               `bson:"rt,omitempty"`
	Output map[string]string `bson:"result,omitempty"`
	Inputs map[string]string `bson:"inputs,omitempty"`
}

func CreateTask(pID primitive.ObjectID, ty string, ns []NodeResult) (primitive.ObjectID, error) {
	t := &Task{}
	t.PipelineId = pID
	t.Resutls = ns
	t.Type = ty
	t.CreateTime = time.Now().In(common.GetLocalZone())
	t.UpdateTime = time.Now().In(common.GetLocalZone())
	return commonCreate(t, common.TaskCollection)
}

func UpdateTask(t *Task) error {
	t.UpdateTime = time.Now().In(common.GetLocalZone())
	return commontUpdate(t, t.ID, common.TaskCollection)
}

func GetTaskByID(id string) (Task, error) {
	ido, err := getIDFromInterface(id)
	if err != nil {
		return Task{}, err
	}
	return commontGetByID[Task](ido, bson.M{}, common.TaskCollection)
}

func GetTasks(q bindmodels.Query) ([]Task, error) {
	filterOption, err := getFilterOption(q)
	if err != nil {
		return []Task{}, err
	}
	oid, err := getIDFromInterface(q.PipelineId)
	if err != nil {
		return nil, err
	}
	filterOption["pipelineId"] = oid
	option := getPagingOption(q)
	option.Projection = bson.M{"points": 0}
	return commonGet[Task](common.TaskCollection, filterOption, option)
}

func GetTaskRunningMonitorInfo(tq bindmodels.TaskQuery) ([]MonitorData, error) {
	oid, err := getIDFromInterface(tq.ID)
	if err != nil {
		return nil, err
	}
	t, err := commontGetByID[Task](oid, bson.M{}, common.TaskCollection)
	if err != nil {
		return nil, err
	}
	pid, err := getIDFromInterface(t.PipelineId)
	if err != nil {
		return nil, err
	}
	p, err := commontGetByID[Pipeline](pid, bson.M{}, common.PipelineCollection)
	if err != nil {
		return nil, err
	}
	if t.Points == nil {
		return []MonitorData{}, nil
	}
	match := &mongoutil.Match{
		Querys: []mongoutil.Query{
			&mongoutil.ComparisonQuery{
				Name: "_id",
				Comparisons: []mongoutil.Comparison{
					{
						Operator: mongoutil.InOperator,
						Value:    t.Points,
					},
				},
			},
		},
	}
	var mds []MonitorData
	for _, n := range p.Flow.Nodes {
		mo, err := getPoints[TypePoint](n.ID, match, tq)
		if err != nil {
			return nil, err
		}
		if len(mo) == 0 {
			return nil, errors.New("get points error")
		}
		mds = append(mds, MonitorData{
			ID:   n.ID,
			Name: n.Data.Name,
			TP:   mo[0],
		})
	}
	return mds, nil
}

func GetTaskRunningFlowInfo(tq bindmodels.TaskQuery) ([]NodeResult, error) {
	oid, err := getIDFromInterface(tq.ID)
	if err != nil {
		return nil, err
	}
	opt := options.FindOneOptions{
		Projection: bson.M{"results": 1, "_id": 1},
	}
	task, err := commontGetByID[Task](oid, bson.M{}, common.TaskCollection, &opt)
	if err != nil {
		return nil, err
	}
	return task.Resutls, nil
}

func GetTaskCount(pipelineId string) (int64, error) {
	ido, err := getIDFromInterface(pipelineId)
	if err != nil {
		return 0, nil
	}
	ctx, cancel, c := getCollection(common.TaskCollection)
	defer cancel()
	count, err := c.CountDocuments(ctx, bson.M{"pipelineId": ido})
	if err != nil {
		return 0, errors.New("get plugin count error")
	}
	return count, nil
}

func UpdatePerformance(id string, points []Point) error {
	ido, err := getIDFromInterface(id)
	if err != nil {
		return err
	}
	ctx, cancel := Mongo.getConnectCtx()
	defer cancel()
	err = Trancaction(ctx, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction()
		if err != nil {
			return err
		}
		pc := Mongo.client.Database(DataBase).Collection(common.PointCollection)

		var ps []interface{}
		for _, p := range points {
			ps = append(ps, p)
		}
		ims, err := pc.InsertMany(ctx, ps)
		if err != nil {
			return err
		}
		tc := Mongo.client.Database(DataBase).Collection(common.TaskCollection)
		_, err = tc.UpdateOne(ctx, bson.M{"_id": ido}, getTaskPointUpdateData(ims.InsertedIDs))

		if err != nil {
			if atErr := sc.AbortTransaction(sc); err != nil {
				err = atErr
			}
			return err
		} else {
			err := sc.CommitTransaction(sc)
			if err != nil {
				return errors.New("commit transaction err: update performance data error")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func getTaskPointUpdateData(ids []interface{}) bson.M {
	return bson.M{
		"$push": bson.M{
			"points": bson.M{
				"$each": ids,
			},
		},
	}
}

func getTStage(id string, tq bindmodels.TaskQuery) map[string][]mongoutil.Stage {
	result := map[string][]mongoutil.Stage{}
	cq := &mongoutil.ComparisonQuery{
		Name:        "timestamp",
		Comparisons: []mongoutil.Comparison{},
	}
	if tq.LastTime > 0 {
		lastTime := int64(tq.LastTime * 1000)
		cs := []mongoutil.Comparison{
			{
				Operator: mongoutil.GteOperator,
				Value:    primitive.DateTime(time.Now().UnixMilli() - lastTime),
			},
		}
		cq.Comparisons = cs
	} else {
		cs := []mongoutil.Comparison{
			{
				Operator: mongoutil.GteOperator,
				Value:    primitive.DateTime(tq.StartTime * 1000),
			},
			{
				Operator: mongoutil.LteOperator,
				Value:    primitive.DateTime(tq.EndTime * 1000),
			},
		}
		cq.Comparisons = cs
	}
	for _, tp := range tps {
		stages := []mongoutil.Stage{
			&mongoutil.Match{
				Querys: []mongoutil.Query{
					&mongoutil.EqualityCondition{
						Key:   "metadata.type",
						Value: tp,
					},
					&mongoutil.EqualityCondition{
						Key:   "metadata.id",
						Value: id,
					},
					cq,
				},
			},
			&mongoutil.Project{
				Specifications: []mongoutil.Specification{
					{
						Key:   "timestamp",
						Value: 1,
					},
					{
						Key:   "temp",
						Value: 1,
					},
				},
			},
			&mongoutil.Sort{
				Fields: map[string]mongoutil.Order{
					"timestamp": {
						IsMeta: false,
						Value:  1,
					},
				},
			},
		}
		result[tp] = stages
	}
	return result
}

func getPoints[T TypePoint](id string, match mongoutil.Stage, tq bindmodels.TaskQuery) ([]T, error) {
	stages := getTStage(id, tq)
	f := &mongoutil.Facet{
		Stages: stages,
	}
	pipeline := []bson.M{match.GetBson(), f.GetBson()}
	rr, err := json.Marshal(pipeline)
	ms, err := commonAggregate[T](common.PointCollection, pipeline)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
