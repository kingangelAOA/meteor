package engine

import (
	"common"
	"context"
	"fmt"
	"sort"
	"time"
)

var SM *StatisticsManager

func init() {
	SM = &StatisticsManager{
		tasks: map[string]*TaskStatistics{},
	}
}

type StatisticsManager struct {
	tasks map[string]*TaskStatistics
}

func (sm *StatisticsManager) AddTask(ctx context.Context, id string, nodeIds []string) {
	sm.tasks[id] = NewTaskStatistics(ctx, nodeIds)
}

func (sm *StatisticsManager) CleanTask(id string) {
	delete(sm.tasks, id)
}

func (sm *StatisticsManager) UpdatePTaskStatistics(taskId, nodeId string, nns NodeStatistics) {
	if _, ok := sm.tasks[taskId]; ok {
		sm.tasks[taskId].AddPNodeStatistics(nodeId, nns)
	}
}

type TaskStatistics struct {
	PNodes    map[string]chan NodeStatistics
	processed map[string]*Processed
	outputs   chan []Point
	ctx       context.Context
}

func NewTaskStatistics(ctx context.Context, nodeIds []string) *TaskStatistics {
	ts := &TaskStatistics{
		PNodes:    map[string]chan NodeStatistics{},
		processed: map[string]*Processed{},
		outputs:   make(chan []Point, common.StatisticsDataLimit),
		ctx:       ctx,
	}
	for _, nodeId := range nodeIds {
		ts.processed[nodeId] = NewProcessed(nodeId)
		ts.PNodes[nodeId] = make(chan NodeStatistics, common.StatisticsDataLimit)
	}
	ts.processing()
	return ts
}

func (ts *TaskStatistics) recieve(ps []Point) {
	if !common.IsClosed(ts.outputs) {
		ts.outputs <- ps
	}
}

func (ts *TaskStatistics) processing() {
	for key, ns := range ts.PNodes {
		go ts.processed[key].Add(ts.ctx, ns)
		go func() {
			<-ts.ctx.Done()
			for k, _ := range ts.PNodes {
				if !common.IsClosed(ts.PNodes[k]) {
					close(ts.PNodes[k])
				}
			}
			if !common.IsClosed(ts.outputs) {
				close(ts.outputs)
			}
		}()
		go ts.processed[key].calculate(ts.ctx)
		go ts.processed[key].output(ts.recieve)
	}
}

func (ts *TaskStatistics) AddPNodeStatistics(nodeId string, nns NodeStatistics) {
	if _, ok := ts.PNodes[nodeId]; ok {
		if !common.IsClosed(ts.PNodes[nodeId]) {
			ts.PNodes[nodeId] <- nns
		}
	} else {
		ts.PNodes[nodeId] = make(chan NodeStatistics, 10000)
		ts.PNodes[nodeId] <- nns
	}
}

type NodeStatistics struct {
	rt     int
	status string
	err    string
	// 毫秒
	startTime int64
}

func NewNodeStatistics(rt int, startTime int64, status, err string) *NodeStatistics {
	return &NodeStatistics{
		rt:        rt,
		status:    status,
		err:       err,
		startTime: startTime,
	}
}

type Processed struct {
	ID     string
	Times  map[int64]int8
	RTS    map[int64]map[int]int
	QPS    map[int64]int
	Err    map[int64]int
	points chan []Point
	lock   chan uint8
}

type Point struct {
	ID        string  `json:"id,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
	Type      string  `json:"type,omitempty"`
	Temp      float64 `json:"temp"`
}

func newPoint(id, t string, timestamp int64, temp float64) Point {
	return Point{
		ID:        id,
		Timestamp: timestamp,
		Type:      t,
		Temp:      temp,
	}
}

func NewProcessed(id string) *Processed {
	lock := make(chan uint8, 1)
	lock <- 0
	return &Processed{
		ID:     id,
		Times:  map[int64]int8{},
		RTS:    map[int64]map[int]int{},
		QPS:    map[int64]int{},
		Err:    map[int64]int{},
		points: make(chan []Point, 10000),
		lock:   lock,
	}
}

func (p *Processed) cleanTime(t int64) {
	delete(p.Times, t)
	delete(p.RTS, t)
	delete(p.QPS, t)
	delete(p.Err, t)
}

func (p *Processed) Add(ctx context.Context, ns chan NodeStatistics) {
	for {
		select {
		case <-ctx.Done():
			return
		case n := <-ns:
			<-p.lock
			st := n.startTime / 1000
			p.Times[st] = 1
			if count, ok := p.QPS[st]; ok {
				p.QPS[st] = count + 1
			} else {
				p.QPS[st] = 1
			}
			if rtMap, ok := p.RTS[st]; ok {
				if count, ok := rtMap[n.rt]; ok {
					p.RTS[st][n.rt] = count + 1
				} else {
					p.RTS[st][n.rt] = 1
				}
			} else {
				p.RTS[st] = map[int]int{}
				p.RTS[st][n.rt] = 1
			}
			if n.status == common.Failed {
				if count, ok := p.Err[st]; ok {
					p.Err[st] = count + 1
				} else {
					p.Err[st] = 1
				}
			}
			p.lock <- 0
		}
	}
}

func (p *Processed) calculate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(p.points)
			return
		default:
			<-p.lock
			now := time.Now().Unix()
			for t, _ := range p.Times {
				if t <= now-common.StatisticsDelay {
					rtCount := 0
					var points []Point
					if rtMap, ok := p.RTS[t]; ok {
						var keys []int
						allRT := 0
						for k, count := range rtMap {
							keys = append(keys, k)
							rtCount = rtCount + count
							allRT = allRT + k*count
						}
						sort.Sort(sort.Reverse(sort.IntSlice(keys)))
						flag99 := int(float64(rtCount) * 0.01)
						flag95 := int(float64(rtCount) * 0.05)
						flag90 := int(float64(rtCount) * 0.10)
						points = append(points, newPoint(p.ID, common.PointTypeRT90, t*1000, float64(getRT(keys, rtMap, flag90))))
						points = append(points, newPoint(p.ID, common.PointTypeRT95, t*1000, float64(getRT(keys, rtMap, flag95))))
						points = append(points, newPoint(p.ID, common.PointTypeRT99, t*1000, float64(getRT(keys, rtMap, flag99))))
						points = append(points, newPoint(p.ID, common.PointTypeRTAverage, t*1000, float64(allRT/rtCount)))
					}
					errRate := 0.0
					if v, ok := p.Err[t]; ok {
//						if rtCount != 0 {
//							errRate = float64(v) / float64(v + rtCount)
//						}
						errRate = float64(v) / float64(rtCount)
						fmt.Println(fmt.Sprintf("t:%d, rtCount: %d, errCount: %d, errRate: %f", t, rtCount, v, errRate))
					}
//					rts, _ := json.Marshal(p.RTS)
//					fmt.Println(fmt.Sprintf("rts: %s", string(rts)))
//					errs, _ := json.Marshal(p.Err)
//					fmt.Println(fmt.Sprintf("errs: %s", string(errs)))
//					fmt.Println("**********************")
					points = append(points, newPoint(p.ID, common.PointTypeErrRate, t*1000, errRate))
					qps := 0
					if v, ok := p.QPS[t]; ok {
						qps = v
					}
					points = append(points, newPoint(p.ID, common.PointTypeQPS, t*1000, float64(qps)))
					p.points <- points
					p.cleanTime(t)
				}
			}
			p.lock <- 0
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func (p *Processed) output(recieve func([]Point)) {
	for {
		select {
		case ps, ok := <-p.points:
			if !ok {
				return
			}
			recieve(ps)
		}
	}
}

func getRT(keys []int, rtMap map[int]int, flagRT int) int {
	flag := 0
	for _, k := range keys {
		flag = flag + rtMap[k]
		if flag >= flagRT {
			return k
		}
	}
	return 0
}
