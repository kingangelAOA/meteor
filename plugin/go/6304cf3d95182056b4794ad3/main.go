package main

import (
	"errors"
	"plugin/shared"

	"github.com/hashicorp/go-plugin"
)

type Http struct {
}

func (h *Http) Run(data map[string]string) (map[string]string, error) {
	return data, errors.New("this is test")
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
