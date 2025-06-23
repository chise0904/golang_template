package trace

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
)

func NewTime() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// SubTimeFromContext get time sub from context
func SubTimeFromContext(ctx context.Context) (milliseconds int64) {
	v, ok := ctx.Value(ctxKeyXTime).(int64)
	if !ok {
		v = NewTime()
		ctx = context.WithValue(ctx, ctxKeyXTime, v)
	}
	milliTime := time.Now().UTC().UnixNano()/1e6 - v
	return milliTime
}

func GetTimeFromContext(ctx context.Context) (milliseconds int64) {
	v, ok := ctx.Value(ctxKeyXTime).(int64)
	if !ok {
		v = NewTime()
	}
	return v
}

// ContextWithTime returns a context.Context with given time value.
func ContextWithTime(ctx context.Context, milliseconds int64) context.Context {
	return context.WithValue(ctx, ctxKeyXTime, milliseconds)
}

// ContextWithTimeForGRPC returns a context.Context with given Time value.
func ContextWithTimeForGRPC(ctx context.Context, milliseconds int64) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metadataTimeKey, strconv.FormatInt(milliseconds, 10))
}

// GetTimeFromContextForGRPC get time sub from meta
func GetTimeFromContextForGRPC(ctx context.Context) (milliseconds int64) {
	var milliTime int64
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NewTime()
	}
	milliTimeMeta, ok := md[metadataTimeKey]
	if !ok || len(milliTimeMeta) == 0 {
		return NewTime()
	} else {
		n, err := strconv.ParseInt(milliTimeMeta[0], 10, 64)
		if err != nil {
			return NewTime()
		}
		milliTime = n
	}
	return milliTime
}

// SubTimeFromContextForGRPC get time sub from meta
func SubTimeFromContextForGRPC(ctx context.Context) (milliseconds int64) {
	var milliTime int64
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		milliTime = NewTime()
	}
	milliTimeMeta, ok := md[metadataTimeKey]
	if !ok || len(milliTimeMeta) == 0 {
		milliTime = NewTime()
	} else {
		n, err := strconv.ParseInt(milliTimeMeta[0], 10, 64)
		if err != nil {
			milliTime = NewTime()
		}
		milliTime = n
	}
	subTime := time.Now().UTC().UnixNano()/1e6 - milliTime

	return subTime
}
