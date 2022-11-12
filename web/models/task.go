package models

type TaskVO struct {
	Index          int            `json:"index,omitempty"`
	ID             string         `json:"id,omitempty"`
	PipelineId     string         `json:"pipelineId,omitempty"`
	NodeResultList []NodeResultVO `json:"nodeResultList,omitempty"`
	Pipeline       PipelineVO     `json:"pipeline,omitempty"`
	Type           string         `json:"type,omitempty"`
	TaskConsume    int            `json:"taskConsume,omitempty"`
	RunningWorks   int            `json:"runningWorks,omitempty"`
	Status         string         `json:"status,omitempty"`
	CreateTime     string         `json:"createTime,omitempty"`
	UpdateTime     string         `json:"updateTime,omitempty"`
}

type NodeResultVO struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Err             string                 `json:"err,omitempty"`
	Status          string                 `json:"status,omitempty"`
	RT              int                    `json:"rt,omitempty"`
	NodePerformance map[string][][]float64 `json:"nodePerformance,omitempty"`
	Output          map[string]string      `json:"output,omitempty"`
	Inputs          map[string]string      `json:"inputs,omitempty"`
}
