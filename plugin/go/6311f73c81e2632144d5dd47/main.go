package main

import (
	"plugin/shared"
    "plugin/proto"
	"github.com/hashicorp/go-plugin"
    // "errors"
)

type Http struct {
}

func (h *Http) Run(data map[string]string) (*proto.OutPut, error) {
    return &proto.OutPut{
        Data: map[string]string{"a":"b"},
        Rt: 1000,
    }, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.GRPCPlugin{Impl: &Http{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}