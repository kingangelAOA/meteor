package core

import (
	"context"
	"encoding/json"
	"fmt"
	"meteor/models"
	"testing"
	"time"
)

var code = `
fmt := import("fmt")
json := import("json")
ctx["num"] = ctx["num"]+2
prints = append(prints, "1111111")
prints = append(prints, "ssdfsdfsd")
// fmt.println(ctx["c"])
prints = append(prints, string(ctx["c"]))
r := json.encode(ctx)
ctx["aaaa"] = string(r)
`

func BenchmarkTengoScript(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ts := newTengoScript(ctx)
	b.ResetTimer()
	b.SetParallelism(20)

	j := "{\"a\": \"b\", \"c\": 2, \"num\": 0}"
	var d map[string]interface{}
	json.Unmarshal([]byte(j), &d)
	s := &Shared{
		Data: d,
	}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ns := s.CopyShared()
			tm := AcquireScriptMessage("test", TengoType, CopyMap(ns.Data))
			tm.SetData(ns.Data)
			ts.Execute(tm)
			<-tm.Ok
			ns.SetData(tm.GetData())
			// fmt.Println("prints: ", tm.GetPrints())
			// fmt.Println("data: ", ns.Data)
			ReleaseScriptMessage(tm)
			// ReleaseShared(cs)
		}
	})
	b.StopTimer()
}

func BenchmarkScriptServer(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ss, _ := NewScriptService(map[string][]models.BaseScript{TengoType: []models.BaseScript{models.BaseScript{
		ID:   "xx",
		Name: "test",
		Code: code,
	}}}, 1000, ctx)
	ss.Run()
	b.SetParallelism(20)
	s := &Shared{
		Data: map[string]interface{}{"test": "1", "num": 0, "c": 2},
	}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ns := s.CopyShared()
			tm := AcquireScriptMessage("test_xx", TengoType, CopyMap(ns.Data))
			ss.PutMessage(tm)
			ok := <-tm.Ok
			if !ok {
				fmt.Println("error: ", tm.ErrMsg)
			}
			ns.SetData(tm.GetData())
			// fmt.Println("prints: ", tm.GetPrints())
			ReleaseScriptMessage(tm)
		}
	})
}

func newTengoScript(ctx context.Context) *TengoScript {
	ts := NewTengoScript(ctx)
	ts.AddScript("test", code)
	return ts
}
