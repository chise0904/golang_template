package interceptor

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"go.starlark.net/lib/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerRecoveryInterceptor returns a new unary server recovery for panic recovery.
func UnaryServerRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		defer func() {
			if r := recover(); r != nil {
				var msg string
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg = msg + fmt.Sprintf("%s:%d\n", file, line)
				}
				log.Error().Msgf("Opps!! panic: %+v\n%s", r, msg)
				err = status.Errorf(codes.Internal, "%+v", r)
				resp = new(proto.Message)
			}
		}()
		resp, err = handler(ctx, req)

		return resp, err
	}

}

// StreamServerRecoveryInterceptor returns a new stream server interceptor for panic recovery
func StreamServerRecoveryInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		defer func() {
			if r := recover(); r != nil {
				var msg string
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg = msg + fmt.Sprintf("%s:%d\n", file, line)
				}
				log.Error().Msgf("Opps!! panic: %+v\n%s", r, msg)
				err = status.Errorf(codes.Internal, "%v", r)
			}
		}()
		err = handler(srv, stream)

		return err
	}
}
