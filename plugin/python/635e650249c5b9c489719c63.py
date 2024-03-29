#!/usr/bin/env python3
from concurrent import futures
import sys
import time

import grpc

import plugin_pb2
import plugin_pb2_grpc
import json
from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc


class PluginServicer(plugin_pb2_grpc.PluginServicer):
    """Implementation of KV service."""

    def Run(self, request, context):
        output = plugin_pb2.OutPut()
        output.data["a"] = "xxxxxxxxxxxxsssss"

        # context.set_code(grpc.StatusCode.UNKNOWN)
        # context.set_details(json.dumps(dict(output)))

        return output


def serve():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set(
        "plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    plugin_pb2_grpc.add_PluginServicer_to_server(PluginServicer(), server)
    health_pb2_grpc.add_HealthServicer_to_server(health, server)
    server.add_insecure_port('127.0.0.1:12345')
    server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:12345|grpc")
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
