package service

import (
	"context"

	"github.com/chise0904/golang_template/pkg/pagination"
	"github.com/chise0904/golang_template/proto/pkg/identity"
	protoIdentity "github.com/chise0904/golang_template/proto/pkg/identity"
	"github.com/chise0904/golang_template/repository"
)

type IdentityService interface {
	EmailVerification(ctx context.Context, token string) (result int32, page_url string)
	EmailVerificationByCode(ctx context.Context, email string, code string) (result int32, page_url string)
	ChangeEmail(ctx context.Context, id string, email string) (err error)
	ChangePhone(ctx context.Context, id string, phone string) (err error)

	SetPasswordByPassword(ctx context.Context, by string, in *SetPassword) (err error)
	SetPasswordByVerifyCode(ctx context.Context, by string, in *SetPassword) (err error)

	CreateAccount(ctx context.Context, req *identity.CreateAccountRequest) (result *repository.AccountInfo, err error)
	RegisterAccount(ctx context.Context, email string, password string, username string, acType protoIdentity.AccountType, permission *protoIdentity.Permission) (result *repository.AccountInfo, err error)
	GetOne(ctx context.Context, id string) (result *repository.AccountInfo, err error)
	GetAll(ctx context.Context, opts *repository.GetAccountListOpts) (detail_data []repository.AccountInfo, pg *pagination.Pagination, err error)
	DeleteAccount(ctx context.Context, id string) (err error)

	SetAccountBlockStatus(ctx context.Context, id string, isBlocked bool) error

	GetUserProfile(ctx context.Context, id string) (result *repository.AccountProfile, err error)
	ListProfiles(ctx context.Context, req *protoIdentity.ListProfilesRequest) (result []*repository.AccountProfile, pg *pagination.Pagination, err error)
	UpdateUserProfile(ctx context.Context, in *protoIdentity.UpdateProfileRequest) (err error)
	// auth
	CreateAccessTokenByPassword(ctx context.Context, in *Login) (id string, status protoIdentity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error)
	CreateAccessTokenByVeriCode(ctx context.Context, in *Login) (id string, status protoIdentity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (id, token, refresh_token string, token_expired, refresh_expired int64, err error)
	CheckAccessToken(ctx context.Context, token string) (result *CheckAccessToken, err error)
	SendVerificationToEmail(ctx context.Context, in *SendVerificationCode) (err error)
	SendVerificationToPhone(ctx context.Context, in *SendVerificationCode) (err error)

	// third
	CreateGoogleOAuthURL() string
	GetUserInfoFromGoogle(ctx context.Context, code string) (email string, name string, err error)
	GoogleUserRegisAndLogin(ctx context.Context, email string, name string) (id string, status int32, token, refresh_token string, token_expired, refresh_expired int64, err error)

	//
	Version(ctx context.Context) (result *Version)
}

type NotifyService interface {
	SendMail(ctx context.Context, recipient *Recipient, subj, body string, htmlContent string) (err error)
	SendMailInBatch(ctx context.Context, recipients []*Recipient, subj, body string, htmlContent string) (err error)

	SendLoginVerifyCode(ctx context.Context, recipient *Recipient, code string) error
	SendFindPassword(ctx context.Context, recipient *Recipient, code string) error
	SendRegisOrChangeEmail(ctx context.Context, recipient *Recipient, code string) (etk string, err error)
}

type JWTAuthAService interface {
	GenerateJWT(fields map[string]any, expireMinutes int) (string, error)
	ParseJWT(jwtToken string) (*JWTAuthClaims, error)
}
