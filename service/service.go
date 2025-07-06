package service

import (
	"context"

	"github.com/chise0904/golang_template/proto/pkg/identity"
	protoIdentity "github.com/chise0904/golang_template/proto/pkg/identity"
	"github.com/chise0904/golang_template/repository"
)

type IdentityService interface {
	//
	Version(ctx context.Context) (result *Version)
	CreateAccount(ctx context.Context, req *identity.CreateAccountRequest) (result *repository.AccountInfo, err error)
	RegisterAccount(ctx context.Context, email string, password string, username string, acType protoIdentity.AccountType, permission *protoIdentity.Permission) (result *repository.AccountInfo, err error)
	CreateAccessTokenByPassword(ctx context.Context, in *Login) (id string, status protoIdentity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error)
	CreateAccessTokenByVeriCode(ctx context.Context, in *Login) (id string, status protoIdentity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error)
}
