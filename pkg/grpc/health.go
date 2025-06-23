package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func setupHealthChecker(s *grpc.Server) {
	healthService := newHealthCheck()
	grpc_health_v1.RegisterHealthServer(s, healthService)
}

// HealthCheckHandler handler
type HealthCheckHandler struct {
}

// NewHealthCheck new health check
func newHealthCheck() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// Check for grpc health check
func (h HealthCheckHandler) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch for grpc health watch
func (h HealthCheckHandler) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
