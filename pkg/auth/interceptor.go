package auth

import (
	"context"
	"encoding/json"

	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/rs/zerolog/log"

	// "github.com/chise0904/golang_template/pkg/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	userClaimsMetadataKey = "user_claims_md_key"
)

func GetUserClaimsFormIncomingContext(ctx context.Context) (*UserClaims, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.NewError(errors.ErrorInternalError, "missing incoming metadata in ctx")
	}

	userClaims := &UserClaims{}
	var found bool
	for k, v := range md {
		if k == userClaimsMetadataKey && len(v) > 0 {
			found = true
			err := json.Unmarshal([]byte(v[0]), userClaims)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("user claim unmarshal failed: %v", err.Error())
				return nil, errors.NewError(errors.ErrorInternalError, "user claim unmarshal failed")
			}
			break
		}
	}

	if !found {
		return nil, errors.ErrorPageNotFound()
	}

	return userClaims, nil

}

func UnaryClientAuthInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		userClaims, err := GetUserClaimsForContext(ctx)
		if err == nil {
			b, err := json.Marshal(userClaims)
			if err == nil {
				ctx = metadata.AppendToOutgoingContext(ctx, userClaimsMetadataKey, string(b))
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientAuthInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		userClaims, err := GetUserClaimsForContext(ctx)
		if err == nil {

			b, err := json.Marshal(userClaims)
			if err == nil {
				ctx = metadata.AppendToOutgoingContext(ctx, userClaimsMetadataKey, string(b))
			}
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func UnaryServerAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		userClaims, err := GetUserClaimsFormIncomingContext(ctx)
		if err == nil {
			ctx = UserClaimsWithContext(ctx, userClaims)
		}

		return handler(ctx, req)
	}
}

func StreamServerAuthInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(srv, ss)
		}

		// _ss := interceptor.WrapServerStream(ss)
		// userClaims, err := GetUserClaimsFormIncomingContext(_ss.WrappedContext)
		// if err == nil {
		// 	_ss.WrappedContext = UserClaimsWithContext(_ss.WrappedContext, userClaims)
		// }

		return handler(srv, ss)
	}
}
