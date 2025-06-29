package impl

import (
	"context"
	"time"

	"github.com/chise0904/golang_template/pkg/errors"
	libtime "github.com/chise0904/golang_template/pkg/time"
	"github.com/rs/zerolog/log"

	"github.com/chise0904/golang_template/pkg/auth"
	"github.com/chise0904/golang_template/pkg/uid/xid"
	"github.com/chise0904/golang_template/repository"

	//e "gitlab.com/hsf-cloud/e-commerce/user-mgmt-server/utils/email"
	"github.com/chise0904/golang_template/proto/pkg/common"
	"github.com/chise0904/golang_template/proto/pkg/identity"
	// "gitlab.com/hsf-cloud/proto/pkg/identity"
	// "gitlab.com/hsf-cloud/lib/auth"
	// "gitlab.com/hsf-cloud/lib/recommender/gorse"
	// "gitlab.com/hsf-cloud/lib/uid/xid"
)

// RegisterAccount implements service.RegisterAccount
func (s *svc) RegisterAccount(ctx context.Context, email string, password string, username string, acType identity.AccountType, permission *identity.Permission) (result *repository.AccountInfo, err error) {
	log.Ctx(ctx)
	nowMS := time.Now().UnixMilli()
	mPW := TranPW(password, "md5", "normal")
	// 檢查存在性 by email
	// 檢查開通狀態
	ac, rnf, err := s.repo.GetByEmailForNormal(ctx, email)
	if err != nil {
		if !rnf {
			log.Error().Msgf("check exist by email failed: %v", err.Error())
			return nil, err
		} else {
			err = nil
		}
	} else { // 若帳號info已存在
		if ac.RegisMode == "NORMAL" {
			err = errors.NewError(errors.ErrorConflict, "user email already exist.")
			log.Error().Msgf("user email already exist: %v ", email)
			return nil, err
		}
	}
	var p repository.Permission
	if permission != nil {
		p = repository.Permission{
			CanAccessCrossAccount: permission.CanAccessCrossAccount,
			ProductRead:           permission.CanReadProduct,
			ProductRewrite:        permission.CanModifyProduct,
			OrderRead:             permission.CanReadOrder,
			OrderRewrite:          permission.CanModifyOrder,
			SubscribeEmail:        permission.CanReceiveEmails,
			CoMarketing:           permission.CanParticipateInMarketing,
		}
	} else {
		p = repository.Permission{
			CanAccessCrossAccount: false,
			ProductRead:           true,
			ProductRewrite:        true,
			OrderRead:             true,
			OrderRewrite:          true,
			SubscribeEmail:        false,
			CoMarketing:           false,
		}
	}
	acModel := &repository.AccountInfo{
		AccountType: int32(acType),
		RegisMode:   "NORMAL",
		Status:      int32(identity.AccountStatus_AccountStatus_ENABLED),
		Permission:  p,
		Email:       email,
		CreatedAt:   nowMS,
	}
	acModel.Password = mPW
	// 生成ID
	acXID := xid.NewXIDGenerator().GenUID()
	acModel.ID = acXID
	// 創建新帳號
	err = s.repo.UpsertAccount(ctx, acModel)
	if err != nil {
		log.Error().Msgf(err.Error())
		return nil, err
	}
	// 創建profile
	acProfile := &repository.AccountProfile{
		AccountID: acModel.ID,
		UserName:  username,
		Gender:    "unfilled",
		Birthday: repository.Date{
			Day:   1,
			Month: 1,
			Year:  1900,
		},
		CreatedAt: nowMS,
	}
	err = s.repo.UpsertProfile(ctx, acProfile)
	if err != nil {
		log.Error().Msgf(err.Error())
		return nil, err
	}

	// 發驗證信
	err = s.sendEmailVerificationLinkAndCode(ctx, acModel, acProfile, email, "email_verification")
	if err != nil {
		log.Error().Msgf(err.Error())
		return nil, err
	}

	// noti gorse
	// gorseUser := &gorse.User{
	// 	UserID: acModel.ID,
	// }

	// _, err = s.gorseClient.InsertUser(ctx, gorseUser)
	// if err != nil {
	// 	log.Error().Msgf("gorse InsertUser failed: %v", err.Error())
	// 	return nil, err
	// }

	return acModel, err
}

// // SetPasswordByPassword implements service.IdentityService
// func (s *svc) SetPasswordByPassword(ctx context.Context, by string, in *service.SetPassword) (err error) {
// 	l := log.Ctx(ctx)

// 	userClaims, err := auth.GetUserClaimsForContext(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	ac, _, err := s.repo.GetByID(ctx, userClaims.GetUserID())
// 	if err != nil {
// 		l.Error().Msgf(err.Error())
// 		return err
// 	}

// 	// mOP := TranMD5(in.OldPassword)
// 	mOP := TranPW(in.OldPassword, "md5", "normal")
// 	if mOP != ac.Password {
// 		err = errors.NewError(errors.ErrorPasswordNotCorrect, "does not match the old password")
// 		l.Error().Msgf(err.Error())
// 		return err
// 	}
// 	// mNP := TranMD5(in.Password)
// 	mNP := TranPW(in.Password, "md5", "normal")
// 	ac.Password = mNP
// 	err = s.repo.UpdateAccount(ctx, ac)
// 	if err != nil {
// 		return err
// 	}

// 	// noti gorse
// 	gorseUserPatch := &gorse.UserPatch{}

// 	_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
// 	if err != nil {
// 		l.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 		return err
// 	}

// 	err = s.PurgeUserTokens(ctx, ac.ID)
// 	if err != nil {
// 		l.Error().Msgf("PurgeUserTokens failed: %s", err.Error())
// 		return err
// 	}
// 	l.Info().Msgf("PurgeUserTokens success")
// 	return

// }

// // SetPasswordByVerifyCode implements service.IdentityService
// func (s *svc) SetPasswordByVerifyCode(ctx context.Context, by string, in *service.SetPassword) (err error) {
// 	// [X] set by email,phone code
// 	l := log.Ctx(ctx)
// 	var ac *repository.AccountInfo
// 	nowMS := time.Now().UnixMilli()

// 	switch by {
// 	case "email":
// 		ac, _, err = s.repo.GetByEmailForNormal(ctx, in.Email)
// 		if err != nil {
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		evc, _, err := s.repo.GetVeriCodeByActionAndID(ctx, ac.ID, "ev_setpassword")
// 		if err != nil {
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		if in.Code != evc.Code {
// 			err = errors.NewError(errors.ErrorNotAllow, "does not match the ev code")
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		if nowMS-evc.CreatedAt >= int64(s.notifyConfig.Setting.CodeExpiresIn*int(time.Minute.Milliseconds())) {
// 			err = errors.NewError(errors.ErrorNotAllow, "ev code expired")
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		// mNP := TranMD5(in.Password)
// 		mNP := TranPW(in.Password, "md5", "normal")
// 		ac.Password = mNP

// 		err = s.repo.UpdateAccount(ctx, ac)
// 		if err != nil {
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		err = s.repo.DeleteVeriCodeByActionAndID(ctx, evc.AccountID, evc.Action)
// 		if err != nil {
// 			l.Error().Msgf(err.Error())
// 			return err
// 		}

// 		// noti gorse
// 		gorseUserPatch := &gorse.UserPatch{}

// 		_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
// 		if err != nil {
// 			l.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 			return err
// 		}

// 	case "phone":

// 	}

// 	err = s.PurgeUserTokens(ctx, ac.ID)
// 	if err != nil {
// 		l.Error().Msgf("PurgeUserTokens failed: %s", err.Error())
// 		return err
// 	}
// 	l.Info().Msgf("PurgeUserTokens success")

// 	return
// }

// // DeleteAccount implements service.AccountService
// func (s *svc) DeleteAccount(ctx context.Context, id string) (err error) {

// 	l := log.Ctx(ctx)
// 	err = s.repo.DeleteAccount(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	// noti gorse
// 	_, err = s.gorseClient.DeleteUser(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf("gorse DeleteUser failed: %v", err.Error())
// 		return err
// 	}

// 	err = s.PurgeUserTokens(ctx, id)
// 	if err != nil {
// 		l.Error().Msgf("PurgeUserTokens failed: %s", err.Error())
// 		return err
// 	}
// 	l.Info().Msgf("PurgeUserTokens success")

// 	return
// }

// // GetOne implements service.AccountService
// func (s *svc) GetOne(ctx context.Context, id string) (result *repository.AccountInfo, err error) {

// 	ac, _, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return nil, err
// 	}

// 	result = ac

// 	return
// }

// // GetAll implements service.AccountService
// func (s *svc) GetAll(ctx context.Context, opts *repository.GetAccountListOpts) (detail_data []repository.AccountInfo, pg *pagination.Pagination, err error) {

// 	detail_data, pg, err = s.repo.GetAccountList(ctx, opts)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return nil, pg, err
// 	}

// 	return
// }

// // GetUserProfile implements service.AccountService
// func (s *svc) GetUserProfile(ctx context.Context, id string) (result *repository.AccountProfile, err error) {

// 	result, err = s.repo.GetProfileByID(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return nil, err
// 	}
// 	return
// }

// // UpdateProfile implements service.AccountService
// func (s *svc) UpdateUserProfile(ctx context.Context, in *protoIdentity.UpdateProfileRequest) (err error) {

// 	tx := s.repo.Begin(ctx)
// 	l := log.Ctx(ctx)
// 	r, err := tx.GetProfileByID(ctx, in.AccountId)
// 	if err != nil {
// 		l.Error().Msgf(err.Error())
// 		tx.Rollback()
// 		return err
// 	}

// 	if in.Birthday != nil {
// 		if !(r.Birthday.Year == 1900 && r.Birthday.Month == 1 && r.Birthday.Day == 1) { // 預設值為 1900-01-01
// 			in.Birthday = nil
// 		}
// 	}

// 	if in.Gender != "" {
// 		if r.Gender != "unfilled" && r.Gender != "" {
// 			in.Gender = ""
// 		}
// 	}

// 	pf := &repository.DBProfileOptions{
// 		AccountID:       in.AccountId,
// 		UserName:        in.UserName,
// 		Icon:            in.Icon,
// 		Description:     in.Description,
// 		Gender:          in.Gender,
// 		Birthday:        s.convertProtoBirthdayValueToRepoBirthday(in.Birthday),
// 		Job:             in.Job,
// 		Country:         in.Country,
// 		City:            in.City,
// 		District:        in.District,
// 		ZipCode:         in.ZipCode,
// 		Address:         in.Address,
// 		ShippingAddress: s.convertProtoShippingAddressValueToRepoAddressArray(in.ShippingAddress),
// 		Language:        in.Language,
// 		Phone_noti:      s.convertProtoBoolTypeToRepoProfileBool(&in.PhoneNoti),
// 		Email_noti:      s.convertProtoBoolTypeToRepoProfileBool(&in.EmailNoti),
// 	}
// 	err = tx.UpdateProfile(ctx, pf)
// 	if err != nil {
// 		l.Error().Msgf(err.Error())
// 		tx.Rollback()
// 		return errors.NewError(errors.ErrorInternalError, err.Error())
// 	}

// 	// noti gorse
// 	gorseUserPatch := &gorse.UserPatch{}

// 	_, err = s.gorseClient.UpdateUser(ctx, in.AccountId, gorseUserPatch)
// 	if err != nil {
// 		l.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 		tx.Rollback()
// 		return err
// 	}

// 	tx.Commit()

// 	err = s.PurgeUserTokens(ctx, in.AccountId)
// 	if err != nil {
// 		l.Error().Msgf("PurgeUserTokens failed: %s", err.Error())
// 		return err
// 	}
// 	l.Info().Msgf("PurgeUserTokens success")

// 	return
// }

// // EmailVerification implements service.AccountService
// func (s *svc) EmailVerification(ctx context.Context, token string) (result int32, page_url string) {
// 	// result : 0:OK 1:TokenError(token&code) 2:NoAccount 3:Verified 4:Freezed 5:InternalError
// 	prefix_url := s.config.IdentityWebURL
// 	lang := s.config.Language

// 	var status string
// 	defer func() {
// 		page_url = fmt.Sprintf("%v/emailVerification/%v/%v", prefix_url, lang, status)
// 	}()

// 	claim, err := s.jwtAuthSvc.ParseJWT(token)
// 	if err != nil {
// 		return 1, ""
// 	}

// 	if claim.IsExpired() {
// 		return 1, ""
// 	}

// 	iemail, ok := claim.Fields["Email"]
// 	if !ok {
// 		return 1, ""
// 	}

// 	email, ok := iemail.(string)
// 	if !ok {
// 		return 1, ""
// 	}

// 	ac, rnf, err := s.repo.GetByEmailForNormal(ctx, email)
// 	if err != nil {
// 		if rnf {
// 			return 2, ""
// 		}
// 		return 5, ""
// 	}

// 	// check account status
// 	switch {
// 	case ac.Status == int32(identity.AccountStatus_AccountStatus_ENABLED) && ac.EvStatus:
// 		return 3, ""
// 	case ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED):
// 		return 4, ""
// 	default:
// 		ac.Status = int32(identity.AccountStatus_AccountStatus_ENABLED)
// 		ac.EvStatus = true
// 	}

// 	// check token match
// 	evc, _, err := s.repo.GetVeriCodeByActionAndID(ctx, ac.ID, "email_verification")
// 	if err != nil {
// 		return 5, ""
// 	}

// 	if token != evc.Token {
// 		return 1, ""
// 	}

// 	// update account
// 	err = s.repo.UpdateAccount(ctx, ac)
// 	if err != nil {
// 		return 5, ""
// 	}

// 	// notify gorse
// 	gorseUserPatch := &gorse.UserPatch{}
// 	_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
// 	if err != nil {
// 		log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 		return 5, ""
// 	}
// 	// remove VerifyCode
// 	err = s.repo.DeleteVeriCodeByActionAndID(ctx, ac.ID, "email_verification")
// 	if err != nil {
// 		return 5, ""
// 	}

// 	status = s.config.StatusOK
// 	return 0, fmt.Sprintf("%v/emailVerification/%v/%v", prefix_url, lang, status)
// }

// // EmailVerificationByCode implements service.IdentityService.
// func (s *svc) EmailVerificationByCode(ctx context.Context, email string, code string) (result int32, page_url string) {
// 	// result : 0:OK 1:TokenError(token&code) 2:NoAccount 3:Verified 4:Freezed 5:InternalError
// 	prefix_url := s.config.IdentityWebURL
// 	lang := s.config.Language
// 	// nowMS := time.Now().UnixMilli()
// 	var ac *repository.AccountInfo
// 	var evc *repository.VerificationCode
// 	var status string
// 	var rnf bool // record not found
// 	var gorseUserPatch *gorse.UserPatch
// 	var err error

// 	ac, rnf, err = s.repo.GetByEmailForNormal(ctx, email)
// 	if err != nil {
// 		if rnf {
// 			result = 2
// 			status = s.config.StatusNoAccount
// 			goto REDIRECT
// 		} else {
// 			result = 5
// 			status = s.config.StatusServerInternalError
// 			goto REDIRECT
// 		}
// 	}
// 	// check account status
// 	switch {
// 	case ac.Status == int32(identity.AccountStatus_AccountStatus_ENABLED) && ac.EvStatus:
// 		result = 3
// 		status = s.config.StatusVerified
// 		goto REDIRECT
// 	case ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED):
// 		result = 4
// 		status = s.config.StatusAccountFreezed
// 		goto REDIRECT
// 	default:
// 		ac.Status = int32(identity.AccountStatus_AccountStatus_ENABLED)
// 		ac.EvStatus = true
// 	}
// 	// check code match
// 	evc, rnf, err = s.repo.GetVeriCodeByActionAndID(ctx, ac.ID, "email_verification")
// 	if err != nil {
// 		if !rnf {
// 			result = 5
// 			status = s.config.StatusServerInternalError
// 			goto REDIRECT
// 		}
// 	}

// 	if code != evc.Code {
// 		result = 1
// 		status = s.config.StatusTokenError
// 		goto REDIRECT
// 	}

// 	// if nowMS-evc.CreatedAt >= int64(e.CODE_EXPIRES_IN*int(time.Minute.Milliseconds())) {
// 	// 	result = 1
// 	// 	status = service.STATUS_TOKEN_ERROR
// 	// 	goto REDIRECT
// 	// }

// 	err = s.repo.UpdateAccount(ctx, ac)
// 	if err != nil {
// 		result = 5
// 		status = s.config.StatusServerInternalError
// 		goto REDIRECT
// 	}

// 	// noti gorse
// 	gorseUserPatch = &gorse.UserPatch{}

// 	_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
// 	if err != nil {
// 		log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 		result = 5
// 		status = s.config.StatusServerInternalError
// 		goto REDIRECT
// 	}

// 	// remove VerifyCode
// 	err = s.repo.DeleteVeriCodeByActionAndID(ctx, ac.ID, "email_verification")
// 	if err != nil {
// 		result = 5
// 		status = s.config.StatusServerInternalError
// 		goto REDIRECT
// 	}

// 	result = 0
// 	status = s.config.StatusOK

// REDIRECT:
// 	page_url = fmt.Sprintf("%v/emailVerification/%v/%v", prefix_url, lang, status)

// 	return result, page_url
// }

// // ================
func (s *svc) sendEmailVerificationLinkAndCode(ctx context.Context, ac *repository.AccountInfo, acProfile *repository.AccountProfile, email string, ev_action string) (err error) {
	// check over request and gen verify code
	nowMS := time.Now().UnixMilli()
	evc, rnf, err := s.repo.GetVeriCodeByActionAndID(ctx, ac.ID, ev_action)
	if err != nil && !rnf {
		return err
	}
	if err == nil && evc != nil {
		if nowMS-evc.CreatedAt < 300000 { // ms
			err = errors.NewError(errors.ErrorNotAllow, "ev_code over request")
			return err
		}
	}

	// recipient := &service.Recipient{
	// 	Name:      acProfile.UserName,
	// 	EmailAddr: email,
	// }
	// code := util.GenRandomCode(8)
	// switch ev_action {
	// case "ev_login":
	// 	err = s.notifySvc.SendLoginVerifyCode(ctx, recipient, code)
	// 	if err != nil {
	// 		log.Error().Msgf(err.Error())
	// 		return err
	// 	}

	// 	// fill in the fields
	// 	if evc == nil {
	// 		evc = &repository.VerificationCode{
	// 			AccountID: ac.ID,
	// 			Action:    ev_action,
	// 		}
	// 	}
	// 	evc.Code = code
	// 	evc.CreatedAt = nowMS

	// case "ev_setpassword":
	// 	err := s.notifySvc.SendFindPassword(ctx, recipient, code)
	// 	if err != nil {
	// 		log.Error().Msgf(err.Error())
	// 		return err
	// 	}

	// 	// fill in the fields
	// 	if evc == nil {
	// 		evc = &repository.VerificationCode{
	// 			AccountID: ac.ID,
	// 			Action:    ev_action,
	// 		}
	// 	}
	// 	evc.Code = code
	// 	evc.CreatedAt = nowMS

	// case "email_verification":

	// 	etk, err := s.notifySvc.SendRegisOrChangeEmail(ctx, recipient, code)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// fill in the fields
	// 	if evc == nil {
	// 		evc = &repository.VerificationCode{
	// 			AccountID: ac.ID,
	// 			Action:    ev_action,
	// 		}
	// 	}
	// 	evc.Code = code
	// 	evc.CreatedAt = nowMS
	// 	evc.Token = etk

	// default:
	// 	return errors.NewError(errors.ErrorInvalidInput, "Invaild ev_action type")
	// }

	err = s.repo.UpsertVeriCodeByActionAndID(ctx, evc)
	if err != nil {
		return err
	}

	return
}

// func (s *svc) convertProtoShippingAddressValueToRepoAddressArray(r []*protoIdentity.Adderss) repository.AddressArray {
// 	var result repository.AddressArray
// 	for _, v := range r {

// 		s := repository.Address{
// 			Type:     v.Type,
// 			Country:  v.Country,
// 			City:     v.City,
// 			District: v.District,
// 			ZipCode:  v.ZipCode,
// 			Address:  v.Address,
// 			StoreID:  v.StoreId,
// 		}

// 		result = append(result, s)
// 	}
// 	return result
// }
// func (s *svc) convertProtoBirthdayValueToRepoBirthday(r *protoIdentity.Date) *repository.Date {
// 	var result *repository.Date

//		if r != nil {
//			result = &repository.Date{
//				Day:   r.Day,
//				Month: r.Month,
//				Year:  r.Year,
//			}
//		} else {
//			return nil
//		}
//		return result
//	}
func (s *svc) convertProtoBoolTypeToRepoProfileBool(r *common.BoolType) *bool {

	var t bool
	switch r.Number() {
	case 0:
		return nil
	case 1:
		t = true
		return &t
	case 2:
		t = false
		return &t
	default:
		return nil
	}

}

// func (s *svc) ListProfiles(ctx context.Context, req *protoIdentity.ListProfilesRequest) (result []*repository.AccountProfile, pg *pagination.Pagination, err error) {

// 	return s.repo.ListProfiles(ctx, &repository.ListProfilesFilter{
// 		AccountIDs: req.AccountIds,
// 		Page:       req.Page,
// 		Perpage:    req.Perpage,
// 	})

// }

// func (s *svc) SetAccountBlockStatus(ctx context.Context, id string, isBlocked bool) error {

// 	l := log.Ctx(ctx)
// 	tx := s.repo.Begin(ctx)
// 	accountInfo, _, err := tx.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	var accountStatus identity.AccountStatus
// 	switch accountInfo.Status {
// 	case int32(identity.AccountStatus_AccountStatus_DELETED):
// 		tx.Rollback()
// 		return nil
// 	case int32(identity.AccountStatus_AccountStatus_ENABLED):
// 		if !isBlocked {
// 			tx.Rollback()
// 			return nil
// 		}
// 		accountStatus = identity.AccountStatus_AccountStatus_BLOCKED
// 	case int32(identity.AccountStatus_AccountStatus_BLOCKED):
// 		if isBlocked {
// 			tx.Rollback()
// 			return nil
// 		}
// 		accountStatus = identity.AccountStatus_AccountStatus_ENABLED

// 	default:
// 		l.Error().Msgf("unknown account status: %s", accountInfo.Status)
// 		return errors.NewError(errors.ErrorInternalError, "unknown account status")

// 	}

// 	err = tx.UpdateAccountStatus(ctx, id, int32(accountStatus))
// 	if err != nil {
// 		tx.Rollback()
// 		l.Error().Msgf("UpdateAccountStatus failed: %s", err.Error())
// 		return err
// 	}
// 	tx.Commit()

// 	if isBlocked {
// 		err = s.PurgeUserTokens(ctx, accountInfo.ID)
// 		if err != nil {
// 			l.Error().Msgf("PurgeUserTokens failed: %s", err.Error())
// 			return err
// 		}

// 	}

// 	return nil

// }

// func (s *svc) PurgeUserTokens(ctx context.Context, userID string) error {
// 	l := log.Ctx(ctx)

// 	l.Debug().Msgf("PurgeUserTokens user id:%s", userID)

// 	tokenListKey := fmt.Sprintf("user-tokens:{%s}", userID)
// 	refreshTokenListKey := fmt.Sprintf("user-refresh-tokens:{%s}", userID)

// 	if deleteTokens, err := s.redisClusterClient.LRange(ctx, tokenListKey, 0, -1).Result(); err != nil {
// 		l.Error().Msgf("LRange failed for %s: %s", tokenListKey, err.Error())
// 		return errors.NewError(errors.ErrorInternalError, err.Error())
// 	} else {
// 		l.Debug().Msgf("delete token: %+v", deleteTokens)
// 		for _, v := range deleteTokens {
// 			err = s.redisClusterClient.Del(ctx, v).Err()
// 			if err != nil {
// 				l.Warn().Msgf("delete token failed: %v", err.Error())
// 			}
// 		}

// 	}

// 	if deleteRefreshTokens, err := s.redisClusterClient.LRange(ctx, refreshTokenListKey, 0, -1).Result(); err != nil {
// 		l.Error().Msgf("LRange failed for %s: %s", refreshTokenListKey, err.Error())
// 		return errors.NewError(errors.ErrorInternalError, err.Error())
// 	} else {
// 		l.Debug().Msgf("delete refresh: %+v", deleteRefreshTokens)
// 		for _, v := range deleteRefreshTokens {
// 			err = s.redisClusterClient.Del(ctx, v).Err()
// 			if err != nil {
// 				l.Warn().Msgf("delete refresh token failed: %v", err.Error())
// 			}
// 		}
// 	}

// 	pipe := s.redisClusterClient.Pipeline()
// 	pipe.Del(ctx, tokenListKey, refreshTokenListKey)

// 	if _, err := pipe.Exec(ctx); err != nil {
// 		l.Warn().Msgf("Pipeline execution failed: %s", err.Error())
// 		return nil
// 	}

// 	l.Info().Msgf("Successfully purged tokens for user %s", userID)
// 	return nil
// }

func (s *svc) CreateAccount(ctx context.Context, req *identity.CreateAccountRequest) (result *repository.AccountInfo, err error) {

	userClaims, err := auth.GetUserClaimsForContext(ctx)
	if err != nil {
		return nil, err
	}

	claims := userClaims.GetAll()
	appID, ok := claims["app_id"].(string)
	now := libtime.NowMS()
	l := log.Ctx(ctx)
	if !ok {
		l.Warn().Msgf("missing app id in ctx")
	}

	tx := s.repo.Begin(ctx)

	mPW := TranPW(req.Password, "md5", "normal")

	acID := xid.NewXIDGenerator().GenUID()

	ac := &repository.AccountInfo{
		ID:          acID,
		AppID:       appID,
		AccountType: int32(req.AccountType),
		RegisMode:   req.RegisMode,
		Status:      int32(req.Status),
		EvStatus:    true,
		PvStatus:    true,
		Permission: repository.Permission{
			CanAccessCrossAccount: req.Permission.CanAccessCrossAccount,
			ProductRead:           req.Permission.CanReadProduct,
			ProductRewrite:        req.Permission.CanModifyProduct,
			OrderRead:             req.Permission.CanReadOrder,
			OrderRewrite:          req.Permission.CanModifyOrder,
			SubscribeEmail:        req.Permission.CanReceiveEmails,
			CoMarketing:           req.Permission.CanParticipateInMarketing,
		},
		Password:  mPW,
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = tx.UpsertAccount(ctx, ac)

	if err != nil {
		l.Error().Msgf("UpsertAccount failed: %v", err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.UpsertProfile(ctx, &repository.AccountProfile{
		AccountID: acID,
		UserName:  req.UserName,
		Birthday: repository.Date{
			Day:   1,
			Month: 1,
			Year:  1900,
		},
		Gender:     req.Gender,
		CreatedAt:  now,
		Phone_noti: *s.convertProtoBoolTypeToRepoProfileBool(&req.PhoneNoti),
		Email_noti: *s.convertProtoBoolTypeToRepoProfileBool(&req.EmailNoti),
	})

	if err != nil {
		l.Error().Msgf("UpsertProfile failed: %v", err.Error())
		tx.Rollback()
		return nil, err
	}

	// gorseUser := &gorse.User{
	// 	UserID: acID,
	// }

	// _, err = s.gorseClient.InsertUser(ctx, gorseUser)
	// if err != nil {
	// 	log.Error().Msgf("gorse InsertUser failed: %v", err.Error())
	// 	return nil, err
	// }

	tx.Commit()
	return ac, nil
}
