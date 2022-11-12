package service

import (
	"common"
	"common/bindmodels"
	"engine"
	"errors"
	"web/db"
	"web/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPipelines(q bindmodels.Query) ([]models.PipelineVO, error) {
	var pvs []models.PipelineVO
	ps, err := db.GetPipelines(q)
	if err != nil {
		return pvs, nil
	}
	for i := 0; i < len(ps); i++ {
		pv := convertToPipelineVO(ps[i])
		pv.Index = i + 1 + int(q.Offset)
		pvs = append(pvs, pv)
	}
	return pvs, nil
}

func CreatePipeline(pv *models.PipelineVO) (primitive.ObjectID, error) {
	p, err := convertToPipelineDB(pv)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return db.CreatePipeline(p)
}

func UpdatePipeline(pv *models.PipelineVO) error {
	p, err := convertToPipelineDB(pv)
	if err != nil {
		return err
	}
	return db.UpdatePipeline(p)
}

func DeletePipeline(ID string) error {
	objectID, err := db.ConverObjectID(ID)
	if err != nil {
		return err
	}
	if err := db.DeletePipeline(objectID); err != nil {
		return err
	}
	return nil
}

func GetPipelinePagination(q bindmodels.Query) (models.Pagination, error) {
	result := models.Pagination{
		Items: []models.PipelineVO{},
		Total: 0,
	}
	pvos, err := GetPipelines(q)
	if err != nil {
		return result, err
	}
	total, err := db.GetPipelineCount()
	if err != nil {
		return result, err
	}
	result.Items = pvos
	result.Total = total
	return result, nil
}

func UpdateFlow(fv *models.FlowVO) error {
	f, err := convertToFlowDB(fv)
	if err != nil {
		return err
	}
	return db.UpdateFlow(f)
}

func GetFlowByPipelineId(pipelineId string) (*models.FlowVO, error) {
	f, err := db.GetFlowByPipelineId(pipelineId)
	if err != nil {
		return nil, err
	}
	if nil == f {
		return nil, nil
	}
	return convertToFlowVO(*f), nil
}

func RunPipeline(pipelineId string) (string, error) {
	pid, err := db.GetIDFromInterface(pipelineId)
	if err != nil {
		return "", err
	}
	p, err := db.GetPipelineByID(pid)
	pv := convertToPipelineVO(p)
	if err != nil {
		return "", err
	}
	fv, err := GetFlowByPipelineId(pipelineId)
	if err != nil {
		return "", err
	}
	if fv == nil {
		return "", errors.New("pipeline has not been set")
	}
	taskId, err := db.CreateTask(pid, pv.Type, getDBNodes(fv))
	if err != nil {
		return "", err
	}
	ns, ls, err := fv.GetNodesLines()
	if err != nil {
		return "", err
	}
	tid := taskId.Hex()
	ep, err := engine.NewPipeline(tid, pv.GetEnginePipelineConfig(), ls, ns)
	if err != nil {
		return "", err
	}
	t, err := engine.NewTask(tid, ep, updateTaskStatistics)
	if err := engine.TM.AddTask(tid, t); err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	if err := t.Run(); err != nil {
		return "", err
	}
	return tid, nil
}

func GetPipelineTypes() []models.PipelineType {
	var types []models.PipelineType
	for k, v := range common.PipelineTypeDec {
		types = append(types, models.PipelineType{
			Name:  v,
			Value: k,
		})
	}
	return types
}

func convertToNodeResultVO(nr db.NodeResult) models.NodeResultVO {
	return models.NodeResultVO{
		ID:     nr.ID,
		Err:    nr.Err,
		RT:     nr.RT,
		Name:   nr.Name,
		Status: nr.Status,
		Output: nr.Output,
		Inputs: nr.Inputs,
	}
}

func convertToFlowVO(f db.Flow) *models.FlowVO {
	var fnvs []models.FlowNodeVO
	for _, n := range f.Nodes {
		fnvs = append(fnvs, convertToFlowNodeVO(n))
	}
	var fevs []models.FlowEdgeVO
	for _, e := range f.Edges {
		fevs = append(fevs, convertToFlowEdgeVO(e))
	}
	return &models.FlowVO{
		PipelineId: f.PipelineId,
		Nodes:      fnvs,
		Edges:      fevs,
	}
}

func convertToFlowEdgeVO(fe db.FlowEdge) models.FlowEdgeVO {
	return models.FlowEdgeVO{
		ID:       fe.ID,
		Type:     fe.Type,
		Source:   fe.Source,
		Target:   fe.Target,
		Animated: fe.Animated,
	}
}

func convertToFlowNodeVO(fn db.FlowNode) models.FlowNodeVO {
	return models.FlowNodeVO{
		ID:       fn.ID,
		Position: fn.Position,
		Type:     fn.Type,
		Data: models.Data{
			DefaultInputs: convertToInputVO(fn.Data.DefaultInputs),
			Describe:      fn.Data.Describe,
			Name:          fn.Data.Name,
			PluginID:      fn.Data.PluginID,
			Code:          fn.Data.Code,
			Language:      fn.Data.Language,
		},
	}
}

func convertToPipelineVO(p db.Pipeline) models.PipelineVO {
	return models.PipelineVO{
		ID:       p.ID.Hex(),
		Describe: p.Describe,
		Type:     p.Type,
		Config: models.Config{
			QPS:   p.Config.QPS,
			Time:  p.Config.Time,
			Users: p.Config.Users,
		},
		CreateTime: p.CreateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
		UpdateTime: p.UpdateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
	}
}

func convertToPipelineDB(pv *models.PipelineVO) (*db.Pipeline, error) {
	p := &db.Pipeline{
		Describe: pv.Describe,
		Type:     pv.Type,
		Config: db.Config{
			QPS:   pv.Config.QPS,
			Users: pv.Config.Users,
			Time:  pv.Config.Time,
		},
	}
	if pv.ID == "" {
		return p, nil
	} else {
		oid, err := db.GetIDFromInterface(pv.ID)
		if err != nil {
			return nil, err
		}
		p.ID = oid
		return p, nil
	}
}

func convertToFlowDB(fv *models.FlowVO) (*db.Flow, error) {
	f := &db.Flow{
		PipelineId: fv.PipelineId,
	}
	var fns []db.FlowNode
	for _, n := range fv.Nodes {
		fns = append(fns, db.FlowNode{
			ID:       n.ID,
			NodeId:   n.NodeId,
			Position: n.Position,
			Type:     n.Type,
			Data: &db.Data{
				DefaultInputs: convertToInputDB(n.Data.DefaultInputs),
				Describe:      n.Data.Describe,
				Name:          n.Data.Name,
				PluginID:      n.Data.PluginID,
			},
		})
	}
	var fes []db.FlowEdge
	for _, e := range fv.Edges {
		fes = append(fes, db.FlowEdge{
			ID:       e.ID,
			Type:     e.Type,
			Source:   e.Source,
			Target:   e.Target,
			Animated: e.Animated,
		})
	}
	f.Nodes = fns
	f.Edges = fes
	return f, nil
}

func convertToInputDB(inputs []models.Input) []db.Input {
	var inputsDB []db.Input
	for _, input := range inputs {
		inputsDB = append(inputsDB, db.Input{
			Name:  input.Name,
			Desc:  input.Desc,
			Value: input.Value,
		})
	}
	return inputsDB
}

func convertToInputVO(inputs []db.Input) []models.Input {
	var inputsVO []models.Input
	for _, input := range inputs {
		inputsVO = append(inputsVO, models.Input{
			Name:  input.Name,
			Desc:  input.Desc,
			Value: input.Value,
		})
	}
	return inputsVO
}

func getDBNodes(fv *models.FlowVO) []db.NodeResult {
	var nrs []db.NodeResult
	for _, n := range fv.Nodes {
		nrs = append(nrs, db.NodeResult{
			ID:   n.ID,
			Name: n.Data.Name,
		})
	}
	return nrs
}

func updateTaskStatistics(ti engine.TaskInfo) error {
	task, err := db.GetTaskByID(ti.ID)
	if err != nil {
		return err
	}
	var results []db.NodeResult
	var ps []db.Point
	isP := false
	for _, ni := range ti.NodeInfoList {
		if len(ni.Points) == 0 {
			results = append(results, db.NodeResult{
				ID:     ni.ID,
				Err:    ni.Err,
				Status: ni.Status,
				RT:     ni.RT,
				Output: ni.Output,
				Inputs: ni.Input,
			})
		} else {
			isP = true
			ps = append(ps, convertToPointDB(ni.Points)...)
		}
	}
	if isP {
		err := UpdateTaskPerformance(ti.ID, ps)
		if err != nil {
			return err
		}
	} else {
		task.Resutls = results
		db.UpdateTask(&task)
	}
	return nil
}

func convertToPointDB(eps []engine.Point) []db.Point {
	var ps []db.Point
	for _, ep := range eps {
		ps = append(ps, db.Point{
			Timestamp: primitive.DateTime(ep.Timestamp),
			Metadata: db.Metadata{
				ID:   ep.ID,
				Type: ep.Type,
			},
			Temp: ep.Temp,
		})
	}
	return ps
}
