package engine

import (
	"common"
	"context"
	"errors"
	"fmt"
)

var TM *TaskManager

func init() {
	TM = &TaskManager{
		tasks: map[string]*Task{},
	}
}

type TaskManager struct {
	tasks map[string]*Task
}

func (tm *TaskManager) GetTaskStatus(id string) string {
	if _, ok := SM.tasks[id]; ok {
		return common.Running
	} else {
		return common.Waiting
	}
}

func (tm *TaskManager) AddTask(id string, task *Task) error {
	if _, ok := tm.tasks[id]; ok {
		return fmt.Errorf("task is already running")
	}
	tm.tasks[id] = task
	return nil
}

func (tm *TaskManager) StopTask(id string) error {
	if task, ok := tm.tasks[id]; ok {
		task.stop()
		cleanSMTM(id)
		return nil
	}
	return errors.New("stop task error, task is not exist")
}

func (tm *TaskManager) cleanTask(id string) {
	delete(tm.tasks, id)
}

func (tm *TaskManager) ResetQPS(id string, qps int) error {
	if task, ok := tm.tasks[id]; ok {
		return task.p.ResetQPS(qps)
	}
	return errors.New("reset qps error, task is not exist")
}

func (tm *TaskManager) GetUsers(id string) (int, error) {
	if task, ok := tm.tasks[id]; ok {
		return task.p.GetUsers()
	}
	return 0, errors.New("get users error, task is not exist")
}

func (tm *TaskManager) ResetUsers(id string, users int) error {
	if task, ok := tm.tasks[id]; ok {
		return task.p.ResetUsers(users)
	}
	return errors.New("reset users error, task is not exist")
}

func (tm *TaskManager) GetTaskRunningWorks(id string) (int, error) {
	if task, ok := tm.tasks[id]; ok {
		return task.p.GetRunningWorks()
	}
	return 0, errors.New("get task running works error, task is not exist")
}

type Task struct {
	id     string
	p      *Pipeline
	uti    func(ti TaskInfo) error
	tCh    chan Transport
	ts     []Transport
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTask(id string, p *Pipeline, uti func(ti TaskInfo) error) (*Task, error) {
	ctx, cancel := p.getContext(id)
	task := &Task{
		id:     id,
		p:      p,
		uti:    uti,
		tCh:    make(chan Transport, common.TaskTransportLimit),
		ctx:    ctx,
		cancel: cancel,
	}
	return task, nil
}

func (t *Task) Run() error {
	if err := t.p.Run(t.ctx, t.receive); err != nil {
		return err
	}
	switch t.p.Config.Type {
	case common.Default:
		ti := TaskInfo{
			ID:           t.id,
			NodeInfoList: []NodeInfo{},
		}
		for _, transport := range t.ts {
			ti.NodeInfoList = append(ti.NodeInfoList, transport.getTaskInfo())
		}
		t.uti(ti)
		t.stop()
		t.p.clear()
	case common.Performance:
		SM.AddTask(t.ctx, t.id, t.getNodeIDs())
		go t.uploadPerformanceData()
		go t.updateTask()
		go t.monitorstatus()
	}
	return nil
}

func (t *Task) monitorstatus() {
	<-t.ctx.Done()
	t.p.clear()
	fmt.Println("monitor status is done")
	if !common.IsClosed(t.tCh) {
		close(t.tCh)
	}
	cleanSMTM(t.id)
}

func (t *Task) stop() {
	t.cancel()
	fmt.Println("task is stoped")
}

func cleanSMTM(id string) {
	TM.cleanTask(id)
	SM.CleanTask(id)
}

func (t *Task) getNodeIDs() []string {
	var ids []string
	for id, _ := range t.p.Nodes {
		ids = append(ids, id)
	}
	return ids
}

func (t *Task) receive(transport Transport) {
	if t.p.Config.Type == common.Performance {
		t.tCh <- transport
	} else {
		t.ts = append(t.ts, transport)
	}
}

func (t *Task) uploadPerformanceData() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case transport := <-t.tCh:
			t.upload(transport)
		}
	}
}
func (t *Task) upload(transport Transport) {
	SM.UpdatePTaskStatistics(t.id, transport.ID, *transport.getNodeStatistics())
}

func (t *Task) updateTask() {
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			ni := NodeInfo{}
			output := <-SM.tasks[t.id].outputs

			ni.Points = output
			
			t.uti(TaskInfo{
				ID:           t.id,
				NodeInfoList: []NodeInfo{ni},
			})
		}
	}
}

type TaskInfo struct {
	ID           string
	NodeInfoList []NodeInfo
}

type NodeInfo struct {
	ID     string
	Err    string
	Status string
	RT     int
	Points []Point
	Output map[string]string
	Input  map[string]string
}
