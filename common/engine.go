package common

const (
	ErrMsgThreshold = 1000
	ExpiryDuration  = 1
)

const (
	LimitConstantMode = "LIMIT_CONSTANT_MODE"
)

const (
	Success = "success"
	Waiting = "waiting"
	Running = "running"
	Failed  = "failed"
)

const (
	Plugin = "plugin"
)

const (
	Default     = "default"
	Performance = "performance"
)

var PipelineTypeDec = map[string]string{
	Default:     "普通任务",
	Performance: "性能测试",
}

const (
	StatisticsDelay = 10
)

const (
	TaskTransportLimit = 1000
)

const (
	PluginTimeout = 10
)

const (
	StatisticsDataLimit = 100
)

const (
	Trace = "trace"
)
