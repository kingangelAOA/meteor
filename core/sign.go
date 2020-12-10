package core

import (
	"github.com/pkg/errors"
	"time"
)

var (
	TypeNotSurpport = errors.New("")
)

const (
	NO_LIMIT = iota
	LIMIT
)

type item struct {
	status bool
	time   time.Time
}

func (i *item) getStatus() bool {
	return i.status
}

func (i *item) getTime() time.Time {
	return i.time
}

type Sign struct {
	Type         int8    `json:"type"`
	Num          float64 `json:"num"`
	itemCh       chan item
	ready        bool
	readyLock    chan int8
	err          error
}

func (s *Sign) Run() {
	s.init()
	s.generatedSign()
}

func (s *Sign) init() {
	s.itemCh = make(chan item, 1)
	s.readyLock = make(chan int8, 1)
	s.readyLock <- 1
}

func (s *Sign) generatedSign() {
	switch s.Type {
	case LIMIT:
		s.limitSign()
	case NO_LIMIT:
	default:
		s.limitSign()
	}
}

func (s *Sign) limitSign() {
	go func(s *Sign) {
		for {
			s.itemCh <- item{
				status: true,
				time:   time.Now(),
			}
			time.Sleep(s.getSleepTime())
		}
	}(s)
}

func (s *Sign) setReady() {
	<-s.readyLock
	s.ready = true
}

func (s *Sign) Get() bool {
	if !s.ready {
		s.setReady()
	}
	if s.Type == NO_LIMIT {
		return true
	}
	i := <-s.itemCh
	return i.getStatus()
}

func (s *Sign) getSleepTime() time.Duration {
	sl := int(1000000000 / s.Num)
	return time.Duration(sl) * time.Nanosecond
}
