# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import plugin_pb2 as plugin__pb2


class PluginStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Run = channel.unary_unary(
                '/proto.Plugin/Run',
                request_serializer=plugin__pb2.Input.SerializeToString,
                response_deserializer=plugin__pb2.OutPut.FromString,
                )


class PluginServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Run(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_PluginServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Run': grpc.unary_unary_rpc_method_handler(
                    servicer.Run,
                    request_deserializer=plugin__pb2.Input.FromString,
                    response_serializer=plugin__pb2.OutPut.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'proto.Plugin', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Plugin(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Run(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.Plugin/Run',
            plugin__pb2.Input.SerializeToString,
            plugin__pb2.OutPut.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)