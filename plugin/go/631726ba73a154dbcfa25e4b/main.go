package main

import (
	"plugin/proto"
	"plugin/shared"
	"time"

	"github.com/hashicorp/go-plugin"
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

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"plugin": &shared.GRPCPlugin{Impl: &Http{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
