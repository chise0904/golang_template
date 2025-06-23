package grpc

import (
	"github.com/chise0904/golang_template/pkg/auth"
	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/pkg/grpc/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient new grpc client
func NewClient(host string) (*grpc.ClientConn, error) {

	interceptors := []grpc.UnaryClientInterceptor{
		errors.UnaryClientInterceptor(),
		interceptor.UnaryClientTimeInterceptor(),
		interceptor.UnaryClientXRequestIDInterceptor(),
		auth.UnaryClientAuthInterceptor(),
	}

	streamInterceptors := []grpc.StreamClientInterceptor{
		errors.StreamClientInterceptor(),
		interceptor.StreamClientTimeInterceptor(),
		interceptor.StreamClientXRequestIDInterceptor(),
		auth.StreamClientAuthInterceptor(),
	}

	options := grpc.WithChainUnaryInterceptor(interceptors...)
	streamOptions := grpc.WithChainStreamInterceptor(streamInterceptors...)

	// No WithInsecure grpc conn will return error
	conn, err := grpc.Dial(host, options, streamOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}
