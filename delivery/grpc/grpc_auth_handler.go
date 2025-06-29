package grpc

// import (
// 	"context"

// 	"github.com/chise0904/golang_template/service"
// 	"github.com/chise0904/golang_template/pkg/errors"
// 	"github.com/chise0904/golang_template/proto/pkg/identity"
// 	protoIdentity "gitlab.com/hsf-cloud/proto/pkg/identity"
// 	"google.golang.org/protobuf/types/known/emptypb"
// )

// // CheckAccessToken implements identity.IdentityServiceServer.
// func (i *handler) CheckAccessToken(c context.Context, req *protoIdentity.CheckAccessTokenRequest) (resp *identity.CheckAccessTokenResponse, err error) {
// 	r, err := i.svc.CheckAccessToken(c, req.AccessToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &identity.CheckAccessTokenResponse{
// 		AccountId:   r.AccountID,
// 		AppId:       r.AppID,
// 		AccountType: r.AccountType,
// 		RegisMode:   r.RegisMode,
// 		Status:      r.Status,
// 		EvStatus:    r.EvStatus,
// 		PvStatus:    r.PvStatus,
// 		Permission: &identity.Permission{
// 			CanAccessCrossAccount:     r.Permission.CanAccessCrossAccount,
// 			CanReadProduct:            r.Permission.ProductRead,
// 			CanModifyProduct:          r.Permission.ProductRewrite,
// 			CanReadOrder:              r.Permission.OrderRead,
// 			CanModifyOrder:            r.Permission.OrderRewrite,
// 			CanReceiveEmails:          r.Permission.SubscribeEmail,
// 			CanParticipateInMarketing: r.Permission.CoMarketing,
// 		},
// 		Email:     r.Email,
// 		Phone:     r.Phone,
// 		EmailNoti: r.Email_noti,
// 		PhoneNoti: r.Phone_noti,
// 	}, nil
// }

// // LoginAccount implements identity.IdentityServiceServer.
// func (i *handler) LoginAccount(c context.Context, req *protoIdentity.LoginAccountRequest) (resp *protoIdentity.LoginAccountResponse, err error) {
// 	if req.Connection != "email" && req.Connection != "sms" {
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "Invaild connection type")
// 	}
// 	in := &service.Login{
// 		Connection: req.Connection,
// 		Email:      req.Email,
// 		Phone:      req.Phone,
// 		Code:       req.Code,
// 		Password:   req.Password,
// 	}
// 	switch req.GrantType {
// 	case "password":
// 		id, status, acTk, rfTk, tkExp, rfExp, err := i.svc.CreateAccessTokenByPassword(c, in)
// 		if err != nil {
// 			return nil, err
// 		}

// 		resp = &protoIdentity.LoginAccountResponse{
// 			AccountId:       id,
// 			Status:          status,
// 			AccessToken:     acTk,
// 			RefreshToken:    rfTk,
// 			TokenType:       "Bearer",
// 			TokenExpireIn:   tkExp,
// 			RefreshExpireIn: rfExp,
// 		}

// 	case "authorization_code":
// 		id, status, acTk, rfTk, tkExp, rfExp, err := i.svc.CreateAccessTokenByVeriCode(c, in)
// 		if err != nil {
// 			return nil, err
// 		}

// 		resp = &protoIdentity.LoginAccountResponse{
// 			AccountId:       id,
// 			Status:          status,
// 			AccessToken:     acTk,
// 			RefreshToken:    rfTk,
// 			TokenType:       "Bearer",
// 			TokenExpireIn:   tkExp,
// 			RefreshExpireIn: rfExp,
// 		}

// 	case "refresh_token":
// 	default:
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "Invaild grant type")
// 	}

// 	return
// }

// // SendVerificationCode implements identity.IdentityServiceServer.
// func (i *handler) SendVerificationCode(c context.Context, req *protoIdentity.SendVerificationCodeRequest) (resp *emptypb.Empty, err error) {
// 	switch req.Connection {
// 	case "email":

// 		in := &service.SendVerificationCode{
// 			Action: req.Action,
// 			Email:  req.Email,
// 		}
// 		err = i.svc.SendVerificationToEmail(c, in)
// 		if err != nil {
// 			return nil, err
// 		}

// 	case "sms":

// 		in := &service.SendVerificationCode{
// 			Action: req.Action,
// 			Phone:  req.Phone,
// 		}
// 		err = i.svc.SendVerificationToPhone(c, in)
// 		if err != nil {
// 			return nil, err
// 		}

// 	default:
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "Invaild connection type")

// 	}

// 	return &emptypb.Empty{}, nil
// }

// // ChangeContacts implements identity.IdentityServiceServer.
// func (i *handler) ChangeContacts(c context.Context, req *protoIdentity.ChangeContactsRequest) (resp *emptypb.Empty, err error) {
// 	switch req.Category {
// 	case "email":

// 		err = i.svc.ChangeEmail(c, req.AccountId, req.Email)
// 		if err != nil {
// 			return nil, err
// 		}

// 	case "phone":
// 		err = i.svc.ChangePhone(c, req.AccountId, req.Phone)
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "Invaild category")
// 	}

// 	return &emptypb.Empty{}, nil
// }

// func (i *handler) ReFreshToken(c context.Context, req *protoIdentity.ReFreshTokenRequest) (resp *protoIdentity.ReFreshTokenResponse, err error) {

// 	if req.RefreshToken == "" {
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "refresh_token is required")
// 	}

// 	accountID, token, refreshToken, tokenExpireIn, refreshExpireIn, err := i.svc.RefreshAccessToken(c, req.RefreshToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &protoIdentity.ReFreshTokenResponse{
// 		AccountId:       accountID,
// 		AccessToken:     token,
// 		TokenType:       "Bearer",
// 		RefreshToken:    refreshToken,
// 		TokenExpireIn:   tokenExpireIn,
// 		RefreshExpireIn: refreshExpireIn,
// 	}, nil
// }
