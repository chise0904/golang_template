package interceptor

import (
	"context"

	"github.com/chise0904/golang_template/trace"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// UnaryServerXRequestIDInterceptor 從 metadata 拿 x-request-id 並放到 ctx 裡
// 如果 metadata 沒有 x-request-id 則產生 request id 並放進 ctx
func UnaryServerXRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		requestID := trace.XRequestIDFromContextForGRPC(ctx)
		logger := log.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)
		ctx = trace.ContextWithXRequestID(ctx, requestID)

		return handler(ctx, req)
	}
}

func StreamServerXRequestIDInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(srv, ss)
		}

		_ss := WrapServerStream(ss)

		requestID := trace.XRequestIDFromContextForGRPC(_ss.WrappedContext)
		logger := log.With().Str("request_id", requestID).Logger()
		ctx := logger.WithContext(_ss.WrappedContext)
		_ss.WrappedContext = trace.ContextWithXRequestID(ctx, requestID)
		return handler(srv, _ss)
	}
}

// UnaryClientXRequestIDInterceptor 將 x-request-id 設定到 grpc metadata
func UnaryClientXRequestIDInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(trace.ContextWithXRequestIDForGRPC(ctx, trace.XRequestIDFromContext(ctx)), method, req, reply, cc, opts...)
	}
}

func StreamClientXRequestIDInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(trace.ContextWithXRequestIDForGRPC(ctx, trace.XRequestIDFromContext(ctx)), desc, cc, method, opts...)
	}
}
