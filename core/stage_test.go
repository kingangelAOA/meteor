package core

import (
	"context"
	"fmt"
	"meteor/models"
	"testing"
	"time"
)

var Ctx context.Context

func Init() {
	parent := context.Background()
	ctx, _ := context.WithTimeout(parent, 1000*time.Second)
	Ctx = ctx
	ss, _ := NewScriptService(getScripts(), 100, Ctx)
	DefaultScriptService = ss
	DefaultScriptService.Run()
}

func TestSingleCoroutine(t *testing.T) {
	Init()
	s := getShared()
	scs, _ := ConnectSingleCoroutineNode(getNodes(), Ctx, s.CopyShared(), false)
	scs.Run()
	fmt.Println(scs.s.Data)
}

func TestMultiCoroutine(t *testing.T) {
	Init()
	s := getShared()
	nodes := []*ScriptNode{NewScript("1", "test", TengoType), NewScript("2", "test", TengoType)}
	scs, _ := ConnectMultiCoroutineNode(LimitConstantMode, getNodes(), 100, 100, Ctx, s.CopyShared())
	scs.run()
	for {
		select {
		case <-Ctx.Done():
			return
		default:
			for _, n := range nodes {
				fmt.Println(scs.se.getQPS(n.id))
			}
			time.Sleep(1 * time.Second)
			// scs.se.getQPS()
		}
	}
}

func getNodes() []Node {
	nodes := []*ScriptNode{NewScript("1", "test", TengoType), NewScript("2", "test", TengoType)}
	ns := []Node{}
	for _, v := range nodes {
		ns = append(ns, v)
	}
	return ns
}

func getShared() Shared {
	return Shared{
		Data: map[string]interface{}{"test": "1", "num": 0, "c": 2},
	}
}

func getScripts() map[string][]models.BaseScript {
	return map[string][]models.BaseScript{
		TengoType: []models.BaseScript{
			models.BaseScript{
				ID:   "1",
				Name: "test",
				Code: `
				fmt := import("fmt")
				json := import("json")
				ctx["num"] = ctx["num"]+2
				r := json.encode(ctx)
				`,
			},
			models.BaseScript{
				ID:   "2",
				Name: "test",
				Code: `
				fmt := import("fmt")
				json := import("json")
				ctx["num"] = ctx["num"]+2
				r := json.encode(ctx)
				ctx["aaaa"] = string(r)
				`,
			},
		},
	}
}

func Test(t *testing.T) {
	fmt.Println(int(1111) / int(10))
}
