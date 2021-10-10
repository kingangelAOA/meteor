package core

import (
	"context"
	"errors"
	"sort"
	"time"
)

const (
	RT99 = 0.99
	RT95 = 0.95
	RT90 = 0.90
)

type Stat struct {
	endTime time.Time
	elapsed int
	ok      bool
}

type QPS struct {
	time  int64
	count int
}

type StatisticsEngine struct {
	data map[string]*Statistics
}

func NewStatisticalEngine() *StatisticsEngine {
	return &StatisticsEngine{
		data: make(map[string]*Statistics),
	}
}

func (se *StatisticsEngine) add(id string, ctx context.Context) {
	s := NewStatistics(ctx)
	s.run()
	se.data[id] = s
}

func (se *StatisticsEngine) pushStat(id string, s Stat) {
	se.data[id].pushStat(s)
}

func (se *StatisticsEngine) getQPS(id string) []QPS {
	return se.data[id].getQPS()
}

func (se *StatisticsEngine) getRT(id string, wl float64) (int, error) {
	return se.data[id].getRT(wl)
}

func (se *StatisticsEngine) getAverageRT(id string) int {
	return se.data[id].getAverageRT()
}

type TimeLine []int64

type Scale []int

type Statistics struct {
	timeLine       TimeLine
	ends           map[int64]int
	countContainer map[int]int64
	scale          Scale
	total          int64
	statCh         chan Stat
	ctx            context.Context
	l              chan int8
}

func NewStatistics(ctx context.Context) *Statistics {
	s := &Statistics{
		ends:           map[int64]int{},
		countContainer: map[int]int64{},
		statCh:         make(chan Stat),
		ctx:            ctx,
		l:              make(chan int8, 1),
	}
	s.l <- 0
	return s
}

func (st *Statistics) lock() func() {
	<-st.l
	return func() {
		st.l <- 0
	}
}

func (st *Statistics) run() {
	go st.preprocess()
}

func (st *Statistics) pushStat(s Stat) {
	defer st.lock()
	st.statCh <- s
}

func (st *Statistics) preprocess() {
	for {
		select {
		case <-st.ctx.Done():
			close(st.statCh)
			return
		case s := <-st.statCh:
			begin := s.endTime.Unix()
			if count, ok := st.ends[begin]; ok {
				count++
			} else {
				st.ends[begin] = 1
				st.timeLine = append(st.timeLine, begin)
				sort.Sort(st.timeLine)
			}
			if v, ok := st.countContainer[s.elapsed]; ok {
				v++
			} else {
				st.countContainer[s.elapsed] = 1
				st.scale = append(st.scale, s.elapsed)
				sort.Sort(st.scale)
			}
			st.total++
		}
	}
}

func (st *Statistics) getQPS() []QPS {
	defer st.lock()
	tpsList := []QPS{}
	for i := 0; i < len(st.timeLine); i++ {
		t := st.timeLine[i]
		if v, ok := st.ends[t]; ok {
			tpsList = append(tpsList, QPS{time: t, count: v})
		}
	}
	return tpsList
}

func (st *Statistics) getRT(wl float64) (int, error) {
	defer st.lock()
	target := int64(float64(st.total) * wl)
	size := len(st.scale)
	var waterCount int64 = 0
	for i := 0; i < size; i++ {
		if v, ok := st.countContainer[st.scale[i]]; ok {
			waterCount = waterCount + v
		}
		if waterCount >= target {
			return st.scale[i], nil
		}
	}
	return 0, errors.New("get TP error")
}

func (st *Statistics) getAverageRT() int {
	defer st.lock()
	totalTime := int64(0)
	for k, v := range st.countContainer {
		totalTime = totalTime + int64(k)*v
	}
	return int(totalTime) / int(st.total)
}

func (tl TimeLine) Len() int {
	return len(tl)
}
func (tl TimeLine) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TimeLine) Less(i, j int) bool {
	return tl[i] < tl[j]
}

func (s Scale) Len() int {
	return len(s)
}
func (s Scale) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Scale) Less(i, j int) bool {
	return s[i] < s[j]
}
