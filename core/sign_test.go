package core

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

func TestSign_Run(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:14000", nil))
	}()
	s := Sign{
		Type: LIMIT,
		Num:  10,
	}
	s.Run()
	for {
		r := s.Get()
		fmt.Println(r)
	}
}
