package core

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)


func test() {
	t := rand.Intn(1000)
	local := time.Now()
	time.Sleep(time.Duration(t) * time.Millisecond)

	fmt.Println("Hello World!*********", t, local)
}

func TestNewPool(t *testing.T) {
	s := Sign{
		Type: LIMIT,
		Num:  10,
	}
	s.Run()
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:14000", nil))
	}()
	rp, _ := NewPool(100)
	rp.Run()

	go func(rp *WrappedPool) {
		for {
			s.Get()
			rp.AddTask(WrappedTask{
				task:   test,
				status: true,
			})
		}
		rp.AddTask(WrappedTask{
			task:   test,
			status: false,
		})
	}(rp)

	time.Sleep(1000 * time.Second)
}
