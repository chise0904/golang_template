package repository

import (
	"context"

	"github.com/chise0904/golang_template/pkg/pagination"
)

type IdentityRepository interface {

	//for tx...
	Begin(context.Context) IdentityRepository
	Commit() error
	Rollback() error

	// account
	UpsertAccount(ctx context.Context, in *AccountInfo) error

	UpdateAccount(ctx context.Context, in *AccountInfo) error
	DeleteAccount(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (out *AccountInfo, record_not_found bool, err error)
	GetByEmailForNormal(ctx context.Context, email string) (out *AccountInfo, record_not_found bool, err error)
	GetByEmailForThirdParty(ctx context.Context, email string, regis_mode string) (out *AccountInfo, record_not_found bool, err error)

	GetSimpAccountList(ctx context.Context, opts *GetAccountListOpts) (out []SimpAccountInfo, page *pagination.Pagination, err error)
	GetAccountList(ctx context.Context, opts *GetAccountListOpts) (out []AccountInfo, page *pagination.Pagination, err error)
	UpdateAccountStatus(ctx context.Context, id string, accountStatus int32) error
	// profile
	UpsertProfile(ctx context.Context, in *AccountProfile) error

	GetProfileByID(ctx context.Context, id string) (out *AccountProfile, err error)
	ListProfiles(ctx context.Context, filter *ListProfilesFilter) (result []*AccountProfile, pg *pagination.Pagination, err error)
	UpdateProfile(ctx context.Context, in *DBProfileOptions) error
	// verification code
	GetVeriCodeByActionAndID(ctx context.Context, id string, action string) (out *VerificationCode, record_not_found bool, err error)
	UpsertVeriCodeByActionAndID(ctx context.Context, in *VerificationCode) error
	DeleteVeriCodeByActionAndID(ctx context.Context, id string, action string) error

	//
	PingDB() (respTimeMilliSec float32, err error)
}
