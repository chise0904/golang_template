package grpc

import (
	"context"
	// "encoding/json"

	// "github.com/chise0904/golang_template/pkg/errors"
	// "github.com/chise0904/golang_template/pkg/web"
	// "github.com/chise0904/golang_template/proto/pkg/common"
	"github.com/chise0904/golang_template/proto/pkg/identity"
	protoIdentity "github.com/chise0904/golang_template/proto/pkg/identity"
	// "github.com/chise0904/golang_template/repository"
	// "github.com/chise0904/golang_template/service"
	// "google.golang.org/protobuf/types/known/emptypb"
)

// // EmailVerification implements identity.IdentityServiceServer.
// func (i *handler) EmailVerification(c context.Context, req *protoIdentity.EmailVerificationRequest) (resp *protoIdentity.VerificationResponse, err error) {
// 	var result int32
// 	if req.Email != "" && req.Code != "" {
// 		result, _ = i.svc.EmailVerificationByCode(c, req.Email, req.Code)
// 	} else {
// 		result, _ = i.svc.EmailVerification(c, req.AccessToken)
// 	}

// 	resp = &protoIdentity.VerificationResponse{
// 		Result: protoIdentity.VerificationResponse_VeriResult(result),
// 	}

// 	return resp, err
// }

// type getAllResponse struct {
// 	Meta     web.ResponsePayLoadMetaData `json:"meta"`
// 	Accounts []repository.AccountInfo    `json:"accounts"`
// }

// type getAllSimpResponse struct {
// 	Meta     web.ResponsePayLoadMetaData  `json:"meta"`
// 	Accounts []repository.SimpAccountInfo `json:"accounts"`
// }

// // GetAllAccount implements identity.IdentityServiceServer.
// func (i *handler) GetAllAccount(c context.Context, req *protoIdentity.GetAllAccountRequest) (resp *protoIdentity.GetAllAccountResponse, err error) {

// 	opts := &repository.GetAccountListOpts{
// 		IDs: req.UserIds,
// 	}
// 	for _, v := range req.AccountTypes {
// 		opts.AccountTypes = append(opts.AccountTypes, int32(v))
// 	}

// 	switch req.Sort {
// 	case "DESC", "desc":
// 		opts.Sort = "desc"
// 	default:
// 		opts.Sort = "asc"
// 	}

// 	if req.Perpage == 0 {
// 		opts.Perpage = 25
// 	} else {
// 		opts.Perpage = req.Perpage
// 	}
// 	if req.Page <= 1 {
// 		opts.Page = 1
// 	} else {
// 		opts.Page = req.Page
// 	}

// 	switch req.By {
// 	case "id", "email", "time":
// 		opts.By = req.By
// 	default:
// 		opts.By = "time"
// 	}

// 	opts.Filter = req.Filter

// 	detail_data, pg, err := i.svc.GetAll(c, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	r := &getAllResponse{
// 		Meta: web.ResponsePayLoadMetaData{
// 			Pagination: pg,
// 		},
// 		Accounts: detail_data,
// 	}

// 	//TODO: do not use json marshal
// 	b, _ := json.Marshal(r)
// 	_ = json.Unmarshal(b, &resp)
// 	return
// }

// // GetOneAccount implements identity.IdentityServiceServer.
// func (i *handler) GetOneAccount(c context.Context, req *protoIdentity.GetOneAccountRequest) (resp *protoIdentity.AccountInfo, err error) {

// 	r, err := i.svc.GetOne(c, req.AccountId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	b, _ := json.Marshal(r)
// 	_ = json.Unmarshal(b, &resp)

// 	return
// }

// // DeleteOneAccount implements identity.IdentityServiceServer.
// func (i *handler) DeleteOneAccount(c context.Context, req *protoIdentity.DeleteOneAccountRequest) (resp *emptypb.Empty, err error) {

// 	err = i.svc.DeleteAccount(c, req.AccountId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &emptypb.Empty{}, nil
// }

// // GetProfile implements identity.IdentityServiceServer.
// func (i *handler) GetProfile(c context.Context, req *protoIdentity.GetProfileRequest) (resp *protoIdentity.UserProfile, err error) {

// 	r, err := i.svc.GetUserProfile(c, req.AccountId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return i.convertProfileToProto(r), nil
// }

// // UpdateProfile implements identity.IdentityServiceServer.
// func (i *handler) UpdateProfile(c context.Context, req *protoIdentity.UpdateProfileRequest) (resp *emptypb.Empty, err error) {

// 	err = i.svc.UpdateUserProfile(c, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &emptypb.Empty{}, nil
// }

// RegisterAccount implements identity.IdentityServiceServer.
func (i *handler) RegisterAccount(c context.Context, req *protoIdentity.RegisterAccountRequest) (resp *protoIdentity.RegisterAccountResponse, err error) {
	r, err := i.svc.RegisterAccount(c, req.Email, req.Password, req.UserName, identity.AccountType_AccountType_USER, nil)
	if err != nil {
		return nil, err
	}

	resp = &protoIdentity.RegisterAccountResponse{
		Email:    r.Email,
		UserName: req.UserName,
		Href:     "",
	}
	return
}

// SetPassword implements identity.IdentityServiceServer.
// func (i *handler) SetPassword(c context.Context, req *protoIdentity.SetPasswordRequest) (resp *emptypb.Empty, err error) {
// 	switch req.By {
// 	case "password":
// 		token := req.AccessToken
// 		in := &service.SetPassword{
// 			AccessToken: token,
// 			OldPassword: req.OldPassword,
// 			Password:    req.Password,
// 		}
// 		err = i.svc.SetPasswordByPassword(c, req.By, in)
// 		if err != nil {
// 			return nil, err
// 		}

// 	case "email", "phone":

// 		in := &service.SetPassword{
// 			Phone:    req.Phone,
// 			Email:    req.Email,
// 			Code:     req.Code,
// 			Password: req.Password,
// 		}
// 		err = i.svc.SetPasswordByVerifyCode(c, req.By, in)
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		return nil, errors.NewError(errors.ErrorInvalidInput, "Invaild set password type")
// 	}

// 	return &emptypb.Empty{}, nil
// }

// func (i *handler) ListProfiles(ctx context.Context, req *protoIdentity.ListProfilesRequest) (*protoIdentity.ListProfilesResponse, error) {

// 	result, pn, err := i.svc.ListProfiles(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data []*protoIdentity.UserProfile
// 	for _, v := range result {
// 		data = append(data, i.convertProfileToProto(v))
// 	}

// 	return &protoIdentity.ListProfilesResponse{
// 		Meta: &common.Pagination{
// 			TotalCount: int32(pn.TotalCount),
// 			TotalPage:  int32(pn.TotalPage),
// 			Page:       int32(pn.Page),
// 			Perpage:    int32(pn.PerPage),
// 		},
// 		Profiles: data,
// 	}, nil
// }

// func (i *handler) convertProfileToProto(r *repository.AccountProfile) *protoIdentity.UserProfile {
// 	protoProfile := &protoIdentity.UserProfile{
// 		AccountId:   r.AccountID,
// 		UserName:    r.UserName,
// 		Icon:        r.Icon,
// 		Description: r.Description,
// 		Gender:      r.Gender,
// 		Birthday: &protoIdentity.Date{
// 			Day:   r.Birthday.Day,
// 			Month: r.Birthday.Month,
// 			Year:  r.Birthday.Year,
// 		},
// 		Job:      r.Job,
// 		Country:  r.Country,
// 		City:     r.City,
// 		District: r.District,
// 		ZipCode:  r.ZipCode,
// 		Address:  r.Address,

// 		Language:  r.Language,
// 		EmailNoti: r.Email_noti,
// 		PhoneNoti: r.Phone_noti,
// 		CreatedAt: r.CreatedAt,
// 		UpdatedAt: r.UpdatedAt,
// 		DeletedAt: int64(r.DeletedAt),
// 	}

// 	for _, v := range r.ShippingAddress {
// 		protoProfile.ShippingAddress = append(protoProfile.ShippingAddress, &protoIdentity.Adderss{
// 			Type:     v.Type,
// 			Country:  v.Country,
// 			City:     v.City,
// 			District: v.District,
// 			ZipCode:  v.ZipCode,
// 			Address:  v.Address,
// 			StoreId:  v.StoreID,
// 		})
// 	}

// 	return protoProfile
// }

// func (i *handler) SetAccountBlockStatus(ctx context.Context, req *identity.SetAccountBlockStatusRequest) (resp *emptypb.Empty, err error) {

// 	err = i.svc.SetAccountBlockStatus(ctx, req.AccountId, req.IsBlocked)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &emptypb.Empty{}, nil
// }

func (i *handler) CreateAccount(ctx context.Context, req *identity.CreateAccountRequest) (resp *identity.CreateAccountResponse, err error) {

	ac, err := i.svc.CreateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return &identity.CreateAccountResponse{
		AccountId:   ac.ID,
		AccountType: protoIdentity.AccountType(ac.AccountType),
		Email:       ac.Email,
	}, nil

}

// func (i *handler) ModifyPermissions(ctx context.Context, req *identity.ModifyPermissionsRequest) (r *emptypb.Empty, err error) {
// 	return &emptypb.Empty{}, nil
// }
