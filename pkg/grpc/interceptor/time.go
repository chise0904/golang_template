package interceptor

import (
	"context"

	"github.com/chise0904/golang_template/pkg/trace"
	"google.golang.org/grpc"
)

// UnaryServerXRequestIDInterceptor 從 metadata 拿 x-request-id 並放到 ctx 裡
// 如果 metadata 沒有 x-request-id 則產生 request id 並放進 ctx
func UnaryServerTimeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}
		return handler(trace.ContextWithTime(ctx, trace.GetTimeFromContextForGRPC(ctx)), req)
	}
}

func StreamServerTimeInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(srv, ss)
		}

		_ss := WrapServerStream(ss)
		_ss.WrappedContext = trace.ContextWithTime(_ss.WrappedContext, trace.GetTimeFromContextForGRPC(_ss.WrappedContext))
		return handler(srv, _ss)
	}
}

// UnaryClientXRequestIDInterceptor 將 x-request-id 設定到 grpc metadata
func UnaryClientTimeInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(trace.ContextWithTimeForGRPC(ctx, trace.GetTimeFromContextForGRPC(ctx)), method, req, reply, cc, opts...)
	}
}

func StreamClientTimeInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(trace.ContextWithTimeForGRPC(ctx, trace.GetTimeFromContextForGRPC(ctx)), desc, cc, method, opts...)
	}
}
