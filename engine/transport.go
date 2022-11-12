package engine

import "common"

type Transport struct {
	ID        string
	Err       string
	Output    map[string]string
	Input     map[string]string
	RT        int
	StartTime int64
	Status    string
}

func NewErrTransport(err string, rt int, st int64) Transport {
	return Transport{
		Err:       err,
		Status:    common.Failed,
		RT:        rt,
		StartTime: st,
	}
}

func NewDefaultTransport(input, output map[string]string, rt int, st int64) Transport {
	return Transport{
		Status:    common.Success,
		Output:    output,
		Input:     input,
		RT:        rt,
		StartTime: st,
	}
}

func (t *Transport) getNodeStatistics() *NodeStatistics {
	return &NodeStatistics{
		rt:        t.RT,
		status:    t.Status,
		err:       t.Err,
		startTime: t.StartTime,
	}
}

func (t *Transport) getTaskInfo() NodeInfo {
	return NodeInfo{
		ID:     t.ID,
		Err:    t.Err,
		Status: t.Status,
		RT:     t.RT,
		Output: t.Output,
		Input:  t.Input,
	}
}
