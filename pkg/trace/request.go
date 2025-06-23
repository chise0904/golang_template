package trace

import (
	"context"

	"github.com/chise0904/golang_template/pkg/uid"
	"google.golang.org/grpc/metadata"
)

var guid = uid.NewUIDGenerator(uid.GeneratorEnumUUID)

func NewRequestID() string {
	return guid.GenUID()
}

func XRequestIDFromContext(ctx context.Context) string {
	v, ok := ctx.Value(ctxKeyXRequestID).(string)
	if !ok {
		v = NewRequestID()
	}
	return v
}

func ContextWithXRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKeyXRequestID, requestID)
}

func ContextWithXRequestIDForGRPC(ctx context.Context, requestID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, grpcMetadataXRequestIDKey, requestID)
}

func XRequestIDFromContextForGRPC(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NewRequestID()
	}
	requestIDMeta, ok := md[grpcMetadataXRequestIDKey]
	if !ok || len(requestIDMeta) == 0 {
		return NewRequestID()
	}

	return requestIDMeta[0]
}
