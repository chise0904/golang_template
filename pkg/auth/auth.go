package auth

import (
	"context"

	"github.com/chise0904/golang_template/pkg/errors"
)

type UserClaims struct {
	AccessToken string                 `json:"access_token"`
	UserID      string                 `json:"user_id"`
	AccountType int32                  `json:"account_type"`
	Permission  map[string]bool        `json:"permission"`
	All         map[string]interface{} `json:"all"`
}

type userClaimsCtxKey string

const (
	UserClaimsCtxKey userClaimsCtxKey = "userClaimsCtxKey"
)

func NewUserClaims(accessToken, userID string, accountType int32, permission map[string]bool, others map[string]interface{}) *UserClaims {
	u := &UserClaims{
		AccessToken: accessToken,
		UserID:      userID,
		AccountType: accountType,
		Permission:  permission,
		All: map[string]interface{}{
			"user_id":      userID,
			"access_token": accessToken,
			"account_type": accountType,
		},
	}
	if u.Permission == nil {
		u.Permission = map[string]bool{}
	}

	for k, v := range permission {
		u.All[k] = v
	}
	for k, v := range others {
		u.All[k] = v
	}

	return u
}

func (u *UserClaims) GetUserID() string {
	return u.UserID
}

func (u *UserClaims) GetAccessToken() string {
	return u.AccessToken
}

func (u *UserClaims) GetPermission() map[string]bool {
	return u.Permission
}

func (u *UserClaims) GetAccountType() int32 {
	return u.AccountType
}

func (u *UserClaims) GetAll() map[string]interface{} {
	return u.All
}

func UserClaimsGet[T any](u *UserClaims, key string) (T, bool) {
	value, ok := u.All[key]
	if !ok {
		var zero T
		return zero, false
	}

	castedValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return castedValue, true
}

func GetUserClaimsForContext(ctx context.Context) (*UserClaims, error) {
	u, ok := ctx.Value(UserClaimsCtxKey).(*UserClaims)
	if !ok {
		return nil, errors.NewError(errors.ErrorNotAllow, "missing user claims in context")
	}

	return u, nil
}

func UserClaimsWithContext(parent context.Context, claims *UserClaims) context.Context {
	ctx := context.WithValue(parent, UserClaimsCtxKey, claims)
	return ctx
}
