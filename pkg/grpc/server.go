package grpc

import (
	"context"
	"net"

	"github.com/chise0904/golang_template/pkg/auth"
	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/pkg/grpc/interceptor"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewGrpcServer new grpc server and add default interceptor
func NewGrpcServer(config *Config) (*grpc.Server, net.Listener, error) {
	var (
		interceptors []grpc.UnaryServerInterceptor
		listener     net.Listener
		err          error
	)

	listener, err = net.Listen("tcp", config.Port)
	if err != nil {
		log.Error().Msgf("grpc bind failed %v", err)
		return nil, nil, err
	}

	interceptors = []grpc.UnaryServerInterceptor{
		interceptor.UnaryServerXRequestIDInterceptor(),
		interceptor.UnaryServerLoggingInterceptor(config.RequestDump, config.ResponseDump, config.MaxLogBodySize),
		interceptor.UnaryServerTimeInterceptor(),
		errors.UnaryServerErrorInterceptor(),
		interceptor.UnaryServerRecoveryInterceptor(),
		auth.UnaryServerAuthInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		interceptor.StreamServerXRequestIDInterceptor(),
		interceptor.StreamServerLoggingInterceptor(),
		interceptor.StreamServerTimeInterceptor(),
		errors.StreamServerErrorInterceptor(),
		interceptor.StreamServerRecoveryInterceptor(),
		auth.StreamServerAuthInterceptor(),
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	setupHealthChecker(server)
	return server, listener, err
}

// RunGRPC start grpc server by use uber fx
func RunGrpcService(listener net.Listener, service *grpc.Server, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info().Msgf("starting grpc service listen on %s", listener.Addr().String())
				if err := service.Serve(listener); err != nil {
					log.Error().Msgf("failed to start grpc service: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msgf("stopping grpc service.")
			service.GracefulStop()
			return nil
		},
	})
}
