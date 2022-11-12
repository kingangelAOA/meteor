package main

import (
    "plugin/proto"
	"plugin/shared"
    "time"
)

type Plugin struct {

}

func (p *Plugin) Run(map[string]string) (*proto.OutPut, error) {
    time.Sleep(300 * time.Millisecond)
	return &proto.OutPut{
		Data: data,
	}, nil
}