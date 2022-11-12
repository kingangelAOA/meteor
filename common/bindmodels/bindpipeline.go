package bindmodels

type Query struct {
	Limit      int64  `form:"limit,omitempty"`
	Offset     int64  `form:"offset,omitempty"`
	SortField  string `form:"sortField,omitempty"`
	SortType   int    `form:"sortType,omitempty"`
	Search     string `form:"search,omitempty"`
	PipelineId string `form:"pipelineId,omitempty"`
}
