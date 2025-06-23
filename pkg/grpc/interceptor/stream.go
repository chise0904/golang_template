package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type wrapServerStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

func (w *wrapServerStream) Context() context.Context {
	return w.WrappedContext
}
func WrapServerStream(stream grpc.ServerStream) *wrapServerStream {
	if existing, ok := stream.(*wrapServerStream); ok {
		return existing
	}
	return &wrapServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
