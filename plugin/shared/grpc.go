package shared

import (
	"common"
	"plugin/proto"
	"time"

	"golang.org/x/net/context"
)

// GRPCClient is an implementation of Plugin that talks over RPC.
type GRPCClient struct{ client proto.PluginClient }

func (m *GRPCClient) Run(data map[string]string) (*proto.OutPut, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout*time.Millisecond)
	defer cancel()
	output, err := m.client.Run(ctx, &proto.Input{
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Plugin
	proto.UnimplementedPluginServer
}

func (m *GRPCServer) Run(ctx context.Context, input *proto.Input) (*proto.OutPut, error) {
	result, err := m.Impl.Run(input.Data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
