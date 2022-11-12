package shared

import (
	"common"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
)

var PM *PluginManeger

func init() {
	PM = &PluginManeger{
		plugins: make(map[string]*PluginClientVersion),
	}
}

type PluginManeger struct {
	plugins map[string]*PluginClientVersion
}

func (pm *PluginManeger) ClearAllPlugin() error {
	for _, v := range pm.plugins {
		if err := v.closeAll(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManeger) GetPluginVersionNum(id string) int {
	if pcv, ok := pm.plugins[id]; ok {
		return len(pcv.versions)
	}
	return 0
}

func (pm *PluginManeger) DeleteTask(id, taskId, version string) {
	if _, ok := pm.plugins[id]; ok {
		if _, ok := pm.plugins[id].versions[version]; ok {
			pm.plugins[id].versions[version].deleteTask(taskId)
		}
	}
}

func (pm *PluginManeger) CheckPluginHealth(id string) string {
	if pcv, ok := pm.plugins[id]; ok {
		if v, ok := pcv.versions[pcv.newest]; ok {
			if err := v.Check(); err != nil {
				return common.Offline
			}
			return common.Health
		}
		return common.Exception
	}
	return common.NotRelease
}

func (pm *PluginManeger) UpdatePlugin(id string, pc PluginClient) error {
	if pcv, ok := pm.plugins[id]; ok {
		version := uuid.New().String()
		pcv.newest = version
		pcv.versions[version] = pc
	} else {
		pcv, err := NewPluginClientVersion(pc)
		if err != nil {
			return err
		}
		pm.plugins[id] = pcv
	}
	return nil
}

func (pm *PluginManeger) GetNewestVersion(id string) (string, error) {
	if pcv, ok := pm.plugins[id]; ok {
		return pcv.getNewestVersion(), nil
	}
	return "", fmt.Errorf("failed to get the latest version of the plugin(%s), please confirm whether the plugin exists or it is plugin node", id)
}

func (pm *PluginManeger) AcquirePlugin(id, version, taskId string) (Plugin, error) {
	if pcv, ok := pm.plugins[id]; ok {
		return pcv.acquirePlugin(version, taskId)
	}
	return nil, errors.New("failed to get the plug-in, please check whether the plug-in exists")
}

func (pm *PluginManeger) ReleasePlugin(id, version string, p Plugin) error {
	if pcv, ok := pm.plugins[id]; ok {
		return pcv.releasePlugin(version, p)
	}
	return errors.New("the plugin does not exist, please check whether the plugin exists")
}

type PluginClientVersion struct {
	newest   string
	versions map[string]PluginClient
}

func NewPluginClientVersion(pc PluginClient) (*PluginClientVersion, error) {
	version := uuid.New().String()
	pcv := &PluginClientVersion{
		newest:   version,
		versions: map[string]PluginClient{},
	}
	pcv.versions[version] = pc
	pcv.periodicallyClean()
	return pcv, nil
}

func (pcv *PluginClientVersion) getNewestVersion() string {
	return pcv.newest
}

func (pcv *PluginClientVersion) acquirePlugin(version, taskId string) (Plugin, error) {
	if v, ok := pcv.versions[version]; ok {
		return v.AcquirePlugin(taskId)
	}
	return nil, errors.New("newest version not exist")
}

func (pcv *PluginClientVersion) releasePlugin(version string, p Plugin) error {
	if v, ok := pcv.versions[version]; ok {
		v.ReleasePlugin(p)
		return nil
	}
	return nil
}

func (pcv *PluginClientVersion) periodicallyClean() {
	go func() {
		for {
			keys := make([]string, 0, len(pcv.versions))
			for k := range pcv.versions {
				keys = append(keys, k)
			}
			for _, k := range keys {
				if k != pcv.newest {
					if err := pcv.versions[k].Close(); err != nil {
						fmt.Println(err.Error())
					} else {
						delete(pcv.versions, k)
					}
				}
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
}

func (pcv *PluginClientVersion) closeAll() error {
	for _, v := range pcv.versions {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}

type PluginClient interface {
	AcquirePlugin(string) (Plugin, error)
	ReleasePlugin(Plugin)
	Check() error
	Close() error
	addTask(string)
	deleteTask(string)
}

type basePluginClient struct {
	runningNum int64
	tasks      sync.Map
}

func newBasePluginClient() basePluginClient {
	return basePluginClient{}
}

func (bpc *basePluginClient) addTask(id string) {
	if id != "" {
		bpc.tasks.LoadOrStore(id, true)
	}
}

func (bpc *basePluginClient) deleteTask(id string) {
	bpc.tasks.LoadAndDelete(id)
}

func (bpc *basePluginClient) getTasks() []string {
	var ids []string
	bpc.tasks.Range(func(key, value interface{}) bool {
		ids = append(ids, key.(string))
		return true
	})
	return ids
}

func (bpc *basePluginClient) addOne() {
	atomic.AddInt64(&bpc.runningNum, 1)
}

func (bpc *basePluginClient) deleteOne() {
	atomic.AddInt64(&bpc.runningNum, -1)
}

func (bpc *basePluginClient) getRunningNum() int64 {
	return atomic.LoadInt64(&bpc.runningNum)
}

type pluginClientWrap struct {
	basePluginClient
	plugin Plugin
}

func newPluginClientWrap(plugin Plugin) PluginClient {
	l := make(chan uint8, 1)
	l <- 0
	return &pluginClientWrap{
		plugin:           plugin,
		basePluginClient: newBasePluginClient(),
	}
}

func (pc *pluginClientWrap) AcquirePlugin(taskId string) (Plugin, error) {
	pc.addOne()
	return pc.plugin, nil
}

func (pc *pluginClientWrap) ReleasePlugin(plugin Plugin) {
	pc.deleteOne()
}

func (pc *pluginClientWrap) Check() error {
	return nil
}

func (pc *pluginClientWrap) Close() error {
	l := pc.getRunningNum()
	if l == 0 {
		return nil
	}
	return fmt.Errorf("%d tasks are running", l)
}

type grpcPluginClientWrap struct {
	basePluginClient
	pc    *plugin.Client
	cp    plugin.ClientProtocol
	cache sync.Pool
}

func newGRPCPluginClientWrap(pc interface{}) (PluginClient, error) {
	npc := pc.(*plugin.Client)
	protocol, err := npc.Client()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &grpcPluginClientWrap{
		basePluginClient: newBasePluginClient(),
		pc:               npc,
		cp:               protocol,
	}, nil
}

func (gpc *grpcPluginClientWrap) AcquirePlugin(taskId string) (Plugin, error) {
	gpc.addOne()
	gpc.addTask(taskId)
	v := gpc.cache.Get()
	if v == nil {
		raw, err := gpc.cp.Dispense(DefaultGRPC)
		if err != nil {
			return nil, err
		}
		return raw.(Plugin), nil
	}
	return v.(Plugin), nil
}

func (gpc *grpcPluginClientWrap) ReleasePlugin(p Plugin) {
	gpc.deleteOne()
	gpc.cache.Put(p)
}

func (gpc *grpcPluginClientWrap) Check() error {
	if err := gpc.cp.Ping(); err != nil {
		return errors.New("grpc plugin unline")
	}
	return nil
}

func (gpc *grpcPluginClientWrap) Close() error {
	l := len(gpc.getTasks())
	if 0 == l {
		if err := gpc.cp.Close(); err != nil {
			return err
		}
		gpc.pc.Kill()
		return nil
	}
	return nil
}
