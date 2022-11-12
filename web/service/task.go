package service

import (
	"common"
	"common/bindmodels"
	"engine"
	"web/db"
	"web/models"
)

func GetTasks(q bindmodels.Query) ([]models.TaskVO, error) {
	var tvs []models.TaskVO
	ts, err := db.GetTasks(q)
	if err != nil {
		return tvs, err
	}
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		tv := convertToTaskVO(t)
		tv.TaskConsume = t.GetConsume()
		tv.Index = i + 1 + int(q.Offset)
		tvs = append(tvs, tv)
	}
	return tvs, nil
}

func GetTaskPagination(q bindmodels.Query) (models.Pagination, error) {
	result := models.Pagination{
		Items: []models.PipelineVO{},
		Total: 0,
	}
	pvos, err := GetTasks(q)
	if err != nil {
		return result, err
	}
	total, err := db.GetTaskCount(q.PipelineId)
	if err != nil {
		return result, err
	}
	result.Items = pvos
	result.Total = total
	return result, nil
}

func GetTaskRunningFlowInfo(tq bindmodels.TaskQuery) (models.TaskVO, error) {
	nrs, err := db.GetTaskRunningFlowInfo(tq)
	if err != nil {
		return models.TaskVO{}, err
	}
	var nrvs []models.NodeResultVO
	for _, nr := range nrs {
		nrvs = append(nrvs, convertToNodeResultVO(nr))
	}
	return models.TaskVO{
		ID:             tq.ID,
		NodeResultList: nrvs,
	}, nil
}

func GetTaskRunningMonitorInfo(tq bindmodels.TaskQuery) ([]models.NodeResultVO, error) {
	mds, err := db.GetTaskRunningMonitorInfo(tq)
	if err != nil {
		return nil, err
	}
	var result []models.NodeResultVO
	for _, md := range mds {
		np := map[string][][]float64{}
		for t, ps := range md.TP {
			var pData [][]float64
			for _, p := range ps {
				pData = append(pData, []float64{float64(p.Timestamp.Time().UnixMilli()), p.Temp})
			}
			np[t] = pData
		}
		result = append(result, models.NodeResultVO{
			ID:              md.ID,
			Name:            md.Name,
			NodePerformance: np,
		})
	}
	return result, nil
}

func GetRunningWorks(id string) (int, error) {
	works, err := engine.TM.GetTaskRunningWorks(id)
	if err != nil {
		return 0, err
	}
	return works, nil
}

func UpdateTaskPerformance(id string, points []db.Point) error {
	err := db.UpdatePerformance(id, points)
	if err != nil {
		return err
	}
	return nil
}

func StopTask(id string) error {
	return engine.TM.StopTask(id)
}

func ResetQPS(id string, qps int) error {
	return engine.TM.ResetQPS(id, qps)
}

func ResetUsers(id string, users int) error {
	return engine.TM.ResetUsers(id, users)
}

func GetTaskPipelineUsers(id string) (int, error) {
	return engine.TM.GetUsers(id)
}

func GetTaskStatus(id string) string {
	return engine.TM.GetTaskStatus(id)
}

func convertToTaskVO(t db.Task) models.TaskVO {
	tv := models.TaskVO{
		ID: t.ID.Hex(),
		Pipeline: models.PipelineVO{
			Name: t.Pipeline.Name,
		},
		Type:       t.Type,
		Status:     t.GetStatus(),
		CreateTime: t.CreateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
		UpdateTime: t.UpdateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
	}
	return tv
}
