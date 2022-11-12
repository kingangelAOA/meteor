package main

import (
    "plugin/proto"
    "time"
)

type Plugin struct {

}

func (p *Plugin) Run(data map[string]string) (*proto.OutPut, error) {
    time.Sleep(300 * time.Millisecond)
	return &proto.OutPut{
		Data: data,
	}, nil
}