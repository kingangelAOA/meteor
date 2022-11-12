package main

import (
	"plugin/proto"
	"time"

	// "errors"
)

type Http struct {
}

func (h *Http) Run(data map[string]string) (*proto.OutPut, error) {
	time.Sleep(300 * time.Millisecond)
	return &proto.OutPut{
		Data: data,
	}, nil
}
