package grpc

import (
	"github.com/chise0904/golang_template/pkg/grpc"
	protoIdentity "github.com/chise0904/golang_template/proto/pkg/identity"
	"github.com/chise0904/golang_template/service"
	"go.uber.org/fx"
)

type handler struct {
	protoIdentity.UnimplementedIdentityServiceServer
	svc service.IdentityService
}

func RunGrpcIdentityDelivery(config *grpc.Config, svc service.IdentityService, lifecycle fx.Lifecycle) error {
	grpcServer, l, err := grpc.NewGrpcServer(config)
	if err != nil {
		return err
	}

	h := &handler{
		svc: svc,
	}

	protoIdentity.RegisterIdentityServiceServer(grpcServer, h)

	grpc.RunGrpcService(l, grpcServer, lifecycle)

	return nil

}
