package bindmodels

type TaskQuery struct {
	ID        string `form:"id,omitempty"`
	StartTime int64  `form:"startTime,omitempty"`
	EndTime   int64  `form:"endTime,omitempty"`
	LastTime  int    `form:"lastTime,omitempty"`
}

type TaskOperateOption struct {
	ID    string `json:"id,omitempty"`
	QPS   int    `json:"qps,omitempty"`
	Users int    `json:"users,omitempty"`
}
