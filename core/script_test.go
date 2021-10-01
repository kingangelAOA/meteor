package core

import (
	"context"
	"encoding/json"
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

	j := "{\"a\": \"b\", \"c\": 2, \"num\": 0}"
	var d map[string]interface{}
	json.Unmarshal([]byte(j), &d)
	// fmt.Println(err.Error())
	s := &Shared{
		Data: d,
	}
	// fmt.Println(fmt.Printf("s:%p", s))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {

			// fmt.Println(fmt.Printf("ns:%p", ns))
			tm := AcquireTengoMessage("test", s.CopyShared())
			ts.Execute(tm)
			<-tm.Ok
			// if !ok {
			// 	fmt.Println("errr:", tm.ErrMsg)
			// } else {
			// 	fmt.Println("data:", tm.s.Data)
			// }
			// fmt.Println("prints: ", tm.GetPrints())
			cs := ReleaseTengoMessage(tm)
			ReleaseShared(cs)
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
	s := &Shared{
		Data: map[string]interface{}{"test": "1", "num": 0},
	}
	ns := s.CopyShared()
	tm := AcquireTengoMessage("test", ns)
	ss.PutMessage(tm)
	<-tm.Ok
	// if !ok {
	// 	fmt.Println(tm.ErrMsg)
	// } else {
	// 	fmt.Println("data:", tm.s.Data)
	// }
	// fmt.Println("prints: ", tm.GetPrints())
	nns := ReleaseTengoMessage(tm)
	ReleaseShared(nns)
}

// func BenchmarkScriptServer(b *testing.B) {
// 	parent := context.Background()
// 	ctx, cancel := context.WithTimeout(parent, 1000*time.Second)
// 	defer cancel()
// 	ss, _ := NewScriptService(newTengoScript(ctx), 1000, ctx)
// 	ss.Run()
// 	b.SetParallelism(20)
// 	s := &Shared{
// 		Data: map[string]interface{}{"test": "1"},
// 	}
// 	b.RunParallel(func(p *testing.PB) {
// 		for p.Next() {
// 			tm := AcquireTengoMessage("test", s.CopyShared())
// 			ss.PutMessage(tm)
// 			<-tm.Ok
// 			cs := ReleaseTengoMessage(tm)
// 			ReleaseShared(cs)
// 		}
// 	})
// }

func newTengoScript(ctx context.Context) *TengoScript {
	ts := NewTengoScript(ctx)
	ts.AddScript("test", `
	fmt := import("fmt")
	json := import("json")
	ctx["num"] = ctx["num"]+2
	prints = append(prints, "1111111")
	prints = append(prints, "ssdfsdfsd")
	fmt.println(ctx["c"])
	prints = append(prints, ctx["c"])
	r := json.encode(ctx)
	ctx["aaaa"] = string(r)
	`)
	return ts
}
