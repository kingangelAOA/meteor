package main

import (
	"plugin/proto"
	"time"

	// "errors"
)

type Plugin struct {
}

func (h *Plugin) Run(data map[string]string) (*proto.OutPut, error) {
	time.Sleep(300 * time.Millisecond)
	return &proto.OutPut{
		Data: data,
	}, nil
}

var Plugin Plugin