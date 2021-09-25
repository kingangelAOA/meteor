package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkTengoScript(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ts := newTengoScript(ctx)
	b.ResetTimer()
	b.SetParallelism(20)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tm := AcquireTengoMessage("test", &Shared{
				Data: map[string]interface{}{"test": "1"},
			})
			ts.Execute(tm)
			<-tm.Ok
			ReleaseTengoMessage(tm)
		}
	})
	b.StopTimer()
}

func TestScriptServer(t *testing.T) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ss, _ := NewScriptService(newTengoScript(ctx), 1000, ctx)
	ss.Run()
	tm := AcquireTengoMessage("test", &Shared{
		Data: map[string]interface{}{"test": "1"},
	})
	ss.PutMessage(tm)
	ok := <-tm.Ok
	if !ok {
		fmt.Println(tm.ErrMsg)
	} else {
		fmt.Println(tm.s.Data)
	}
	fmt.Println("prints: ", tm.GetPrints())
	ReleaseTengoMessage(tm)
}

func BenchmarkScriptServer(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ss, _ := NewScriptService(newTengoScript(ctx), 1000, ctx)
	ss.Run()
	b.SetParallelism(20)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tm := AcquireTengoMessage("test", &Shared{
				Data: map[string]interface{}{"test": "1"},
			})
			ss.PutMessage(tm)
			<-tm.Ok
			ReleaseTengoMessage(tm)
		}
	})
}

func newTengoScript(ctx context.Context) *TengoScript {
	ts := NewTengoScript(ctx)
	ts.AddScript("test", `
	json := import("json")
	prints = append(prints, "1111111")
	r := json.encode(ctx)
	ctx["aaaa"] = string(r)
	`)
	return ts
}
