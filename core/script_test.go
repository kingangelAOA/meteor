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
	for n := 0; n < b.N; n++ {
		tm := AcquireTengoMessage("test", &WrappedContext{
			Data: map[string]interface{}{"test": "1"},
		})
		ts.Execute(tm)
		<-tm.Ok
		ReleaseTengoMessage(tm)
		// if !ok {
		// 	fmt.Println(tm.ErrMsg)
		// } else {
		// 	fmt.Println(tm.WrCh.Data)
		// }
		// fmt.Println("prints: ", tm.GetPrings())
	}
	b.StopTimer()
}

func BenchmarkTengoScript1(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ts := newTengoScript(ctx)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tm := &TengoMessage{
			Name: "test",
			WrCh: &WrappedContext{
				Data: map[string]interface{}{"test": "1"},
			},
			Ok: make(chan bool, 1),
		}
		ts.Execute(tm)
		<-tm.Ok
		// if !ok {
		// 	fmt.Println(tm.ErrMsg)
		// } else {
		// 	fmt.Println(tm.WrCh.Data)
		// }
		// fmt.Println("prints: ", tm.GetPrings())
	}
	b.StopTimer()
}

func TestScriptServer(t *testing.T) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ss, _ := NewScriptServer(newTengoScript(ctx), 1000, ctx)
	ss.Run()
	tm := AcquireTengoMessage("test", &WrappedContext{
		Data: map[string]interface{}{"test": "1"},
	})
	ss.PutMessage(tm)
	ok := <-tm.Ok
	if !ok {
		fmt.Println(tm.ErrMsg)
	} else {
		fmt.Println(tm.WrCh.Data)
	}
	fmt.Println("prints: ", tm.GetPrings())
	ReleaseTengoMessage(tm)
}

func BenchmarkScriptServer(b *testing.B) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
	defer cancel()
	ss, _ := NewScriptServer(newTengoScript(ctx), 1000, ctx)
	ss.Run()
	for n := 0; n < b.N; n++ {
		tm := AcquireTengoMessage("test", &WrappedContext{
			Data: map[string]interface{}{"test": "1"},
		})
		ss.PutMessage(tm)
		<-tm.Ok
		ReleaseTengoMessage(tm)
	}

}

func newTengoScript(ctx context.Context) *TengoScript {
	ts := NewTengoScript(ctx)
	ts.AddScript("test", `
	sum := 1 + "2"
	json := import("json")
	prints = append(prints, "1111111")
	r := json.encode(ctx)
	ctx["aaaa"] = string(r)
	`)
	return ts
}
