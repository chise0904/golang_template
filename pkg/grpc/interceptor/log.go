package interceptor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chise0904/golang_template/trace"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// UnaryServerLoggingInterceptor returns a new unary server logging and set request id.
func UnaryServerLoggingInterceptor(reqDump, respDump bool, maxLogBodySize int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip health check
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}
		startTime := time.Now()
		var reqJson []byte
		if reqDump {
			reqJson, _ = json.Marshal(req)
			log.Info().
				Str("method", info.FullMethod).
				Str("request_id", trace.XRequestIDFromContextForGRPC(ctx)).
				Str("grpc.req", string(reqJson)).
				Int64("time_ms", trace.SubTimeFromContextForGRPC(ctx)).
				Msg("Access Log")
		}
		resp, err := handler(ctx, req)
		if respDump && err == nil {
			respJson, _ := json.Marshal(resp)
			if maxLogBodySize > 0 && len(respJson) > maxLogBodySize {
				respJson = append(respJson[:maxLogBodySize], []byte("...")...)
			}
			if reqJson == nil {
				reqJson, _ = json.Marshal(req)
			}

			log.Info().
				Str("grpc.req", string(reqJson)).
				Str("method", info.FullMethod).
				Str("request_id", trace.XRequestIDFromContextForGRPC(ctx)).
				Str("grpc.resp", string(respJson)).
				Float64("grpc.time_ms", float64(time.Since(startTime))/1e6).
				Int64("time_ms", trace.SubTimeFromContextForGRPC(ctx)).
				Msg("Access Log")
		}
		return resp, err
	}
}

func StreamServerLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip health check
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(srv, ss)
		}

		startTime := time.Now()
		_ss := WrapServerStream(ss)

		log.Info().
			Str("method", info.FullMethod).
			Str("request_id", trace.XRequestIDFromContext(_ss.WrappedContext)).
			Int64("time_ms", trace.SubTimeFromContext(_ss.WrappedContext)).
			Msg("Access Log")

		err := handler(srv, ss)

		log.Info().
			Str("method", info.FullMethod).
			Str("request_id", trace.XRequestIDFromContext(_ss.WrappedContext)).
			Float64("grpc.time_ms", float64(time.Since(startTime))/1e6).
			Int64("time_ms", trace.SubTimeFromContext(_ss.WrappedContext)).
			Msg("Access Log")

		return err
	}
}
