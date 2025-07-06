package impl

import (
	"context"
	// "encoding/json"
	"fmt"
	// "strings"
	"time"

	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/repository"
	"github.com/chise0904/golang_template/service"

	// "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"

	// "gitlab.com/hsf-cloud/lib/recommender/gorse"
	libtime "github.com/chise0904/golang_template/pkg/time"
	"github.com/chise0904/golang_template/pkg/uid"
	"github.com/chise0904/golang_template/pkg/uid/uuid"
	"github.com/chise0904/golang_template/proto/pkg/identity"
)

// CreateAccessTokenByPassword implements service.IdentityService
func (s *svc) CreateAccessTokenByPassword(ctx context.Context, in *service.Login) (id string, status identity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error) {

	token_expired = libtime.MilliSecond(time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.AccessTokenExpiresIn)))
	refresh_expired = libtime.MilliSecond(time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)))
	resp := &repository.AccountInfo{}
	switch in.Connection {
	case "email":

		ac, _, err := s.repo.GetByEmailForNormal(ctx, in.Email)
		if err != nil {
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}
		if ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED) {
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorNotAllow, "account is blocked")
		}
		// mPW := TranMD5(in.Password)
		mPW := TranPW(in.Password, "md5", "normal")
		if ac.Password != mPW {
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorPasswordNotCorrect, "password not correct")
		}

		token, refresh_token, err = s.GenToken(ctx, ac)
		if err != nil {
			log.Error().Msgf("gen token failed: %v", err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		}
		// // gen token & refresh token
		// token, err = access_token.GenerateJWT(ac.ID, ac.AccountType, token_expired)
		// if err != nil {
		// 	log.Error().Msgf("gen token failed: %v", err.Error())
		// 	return "", "", "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// }
		// refresh_token, err = access_token.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		// if err != nil {
		// 	log.Error().Msgf("gen token failed: %v", err.Error())
		// 	return "", "", "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// }
		// isExp, err := access_token.IsExpired(ac.RefreshToken)
		// if isExp || err != nil {
		// 	refresh_expired := time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)).UnixMilli()
		// 	ac.RefreshToken, err = access_token.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		// 	if err != nil {
		// 		log.Error().Msgf("gen token failed: %v", err.Error())
		// 		return "", "", "", expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// 	}
		// }

		t := time.Now().UnixMilli()
		ac.LoginAt = t
		ac.UpdatedAt = t

		// update
		err = s.repo.UpdateAccount(ctx, ac)
		if err != nil {
			log.Error().Msgf("update account for login failed: %v", err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}

		resp = ac

		// noti gorse
		// gorseUserPatch := &gorse.UserPatch{}

		// _, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
		// if err != nil {
		// 	log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
		// 	return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		// }

	case "sms":
		return
	}
	return resp.ID, identity.AccountStatus(resp.Status), token, refresh_token, token_expired, refresh_expired, err
}

// // CreateAccessTokenByVeriCode implements service.IdentityService
func (s *svc) CreateAccessTokenByVeriCode(ctx context.Context, in *service.Login) (id string, status identity.AccountStatus, token, refresh_token string, token_expired, refresh_expired int64, err error) {
	token_expired = time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.AccessTokenExpiresIn)).UnixMilli()
	refresh_expired = time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)).UnixMilli()
	nowMS := time.Now().UnixMilli()
	resp := &repository.AccountInfo{}
	switch in.Connection {
	case "email":

		ac, _, err := s.repo.GetByEmailForNormal(ctx, in.Email)
		if err != nil {
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}
		if ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED) {
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorNotAllow, "account is blocked")
		}

		evc, _, err := s.repo.GetVeriCodeByActionAndID(ctx, ac.ID, "ev_login")
		if err != nil {
			log.Error().Msgf(err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}

		if in.Code != evc.Code {
			err = errors.NewError(errors.ErrorNotAllow, "does not match the ev code")
			log.Error().Msgf(err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}

		// if nowMS-evc.CreatedAt >= int64(s.notifyConfig.Setting.CodeExpiresIn*int(time.Minute.Milliseconds())) {
		// 	err = errors.NewError(errors.ErrorNotAllow, "ev code expired")
		// 	log.Error().Msgf(err.Error())
		// 	return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		// }

		token, refresh_token, err = s.GenToken(ctx, ac)
		if err != nil {
			log.Error().Msgf("gen token failed: %v", err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		}

		// // gen token & refresh token
		// token, err = access_token.GenerateJWT(ac.ID, ac.AccountType, token_expired)
		// if err != nil {
		// 	log.Error().Msgf("gen token failed: %v", err.Error())
		// 	return "", "", "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// }
		// refresh_token, err = access_token.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		// if err != nil {
		// 	log.Error().Msgf("gen token failed: %v", err.Error())
		// 	return "", "", "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// }
		// isExp, err := access_token.IsExpired(ac.RefreshToken)
		// if isExp || err != nil {
		// 	refresh_expired := time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)).UnixMilli()
		// 	ac.RefreshToken, err = access_token.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		// 	if err != nil {
		// 		log.Error().Msgf("gen token failed: %v", err.Error())
		// 		return "", "", "", expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// 	}
		// }

		ac.LoginAt = nowMS
		ac.UpdatedAt = nowMS

		// update
		err = s.repo.UpdateAccount(ctx, ac)
		if err != nil {
			log.Error().Msgf("update account for login failed: %v", err.Error())
			return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		}

		resp = ac

		// noti gorse
		// gorseUserPatch := &gorse.UserPatch{}

		// _, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
		// if err != nil {
		// 	log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
		// 	return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		// }

		// err = s.repo.DeleteVeriCodeByActionAndID(ctx, evc.AccountID, evc.Action)
		// if err != nil {
		// 	log.Error().Msgf(err.Error())
		// 	return "", identity.AccountStatus_AccountStatus_UNKNOWN, "", "", token_expired, refresh_expired, err
		// }

	case "sms":
		return
	}
	return resp.ID, identity.AccountStatus(resp.Status), token, refresh_token, token_expired, refresh_expired, err
}

// // SendVerificationToEmail implements service.IdentityService
// func (s *svc) SendVerificationToEmail(ctx context.Context, in *service.SendVerificationCode) (err error) {

// 	ac, _, err := s.repo.GetByEmailForNormal(ctx, in.Email)
// 	if err != nil {
// 		return err
// 	}
// 	if ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED) {
// 		return errors.NewError(errors.ErrorNotAllow, "account is blocked")
// 	}

// 	userProfile, err := s.repo.GetProfileByID(ctx, ac.ID)
// 	if err != nil {
// 		return err
// 	}

// 	switch in.Action {
// 	case "login":

// 		err = s.sendEmailVerificationLinkAndCode(ctx, ac, userProfile, in.Email, "ev_login")
// 		if err != nil {
// 			log.Error().Msgf(err.Error())
// 			return err
// 		}

// 	case "setpassword":

// 		err = s.sendEmailVerificationLinkAndCode(ctx, ac, userProfile, in.Email, "ev_setpassword")
// 		if err != nil {
// 			log.Error().Msgf(err.Error())
// 			return err
// 		}

// 	case "revalidate":
// 		//
// 		ac, _, err := s.repo.GetByEmailForNormal(ctx, in.Email)
// 		if err != nil {
// 			return err
// 		}
// 		if ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED) {
// 			return errors.NewError(errors.ErrorNotAllow, "account is blocked")
// 		}
// 		if ac.EvStatus {
// 			return errors.NewError(errors.ErrorNotAllow, "account already verified")
// 		}
// 		//
// 		err = s.sendEmailVerificationLinkAndCode(ctx, ac, userProfile, in.Email, "email_verification")
// 		if err != nil {
// 			log.Error().Msgf(err.Error())
// 			return err
// 		}
// 	case "auth":
// 	default:
// 		return errors.NewError(errors.ErrorInvalidInput, "Invaild action type")
// 	}
// 	return
// }

// // SendVerificationToPhone implements service.IdentityService
// func (*svc) SendVerificationToPhone(ctx context.Context, in *service.SendVerificationCode) (err error) {
// 	switch in.Action {
// 	case "login":

// 	case "setpassword":

// 	case "auth":

// 	}
// 	return
// }

// // RefreshAccessToken implements service.AuthenticationService
// func (s *svc) RefreshAccessToken(ctx context.Context, refreshToken string) (id, token, refresh_token string, token_expired, refresh_expired int64, err error) {

// 	if strings.Index(refreshToken, "rftk-") < 0 {
// 		return "", "", "", 0, 0, errors.NewError(errors.ErrorNotAllow, "unknown refresh token")
// 	}
// 	claimStr := s.redisClusterClient.Get(ctx, refreshToken).Val()
// 	if claimStr == "" {
// 		return "", "", "", 0, 0, errors.NewError(errors.ErrorNotAllow, "refresh token not found")
// 	}
// 	userClaims := &service.UserClaims{}
// 	err = json.Unmarshal([]byte(claimStr), userClaims)
// 	if err != nil {
// 		return "", "", "", 0, 0, errors.NewError(errors.ErrorNotAllow, err.Error())
// 	}

// 	token_expired = libtime.MilliSecond(time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.AccessTokenExpiresIn)))
// 	refresh_expired = libtime.MilliSecond(time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)))
// 	accessToken, refreshToken, err := s.GenToken(ctx, userClaims.AccountInfo)
// 	if err != nil {
// 		return "", "", "", 0, 0, errors.NewError(errors.ErrorNotAllow, err.Error())
// 	}

// 	return userClaims.AccountInfo.ID, accessToken, refreshToken, token_expired, refresh_expired, nil
// }

// // CheckAccessToken implements service.IdentityService
// func (s *svc) CheckAccessToken(ctx context.Context, token string) (result *service.CheckAccessToken, err error) {

// 	if strings.Index(token, "rftk-") == 0 {
// 		return result, errors.NewError(errors.ErrorNotAllow, "unknown token")
// 	}
// 	tokenSlotKey, err := s.NormalizeTokenSlotKey(token)
// 	if err != nil {
// 		return
// 	}
// 	claimStr := s.redisClusterClient.Get(ctx, tokenSlotKey).Val()
// 	if claimStr == "" {
// 		return nil, errors.NewError(errors.ErrorNotAllow, "token not found")
// 	}

// 	userClaims := &service.UserClaims{}
// 	err = json.Unmarshal([]byte(claimStr), userClaims)
// 	if err != nil {
// 		return nil, errors.NewError(errors.ErrorNotAllow, err.Error())
// 	}

// 	// _, _, err = access_token.GetExpired(token)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// tkID, err := access_token.GetID(token)
// 	// if err != nil {
// 	// 	log.Error().Msgf(err.Error())
// 	// 	return nil, err
// 	// }

// 	// ac, _, err := s.repo.GetByID(ctx, tkID)
// 	// if err != nil {
// 	// 	log.Error().Msgf(err.Error())
// 	// 	return nil, err
// 	// }
// 	// // if ac.AccessToken != token {
// 	// // 	err = errors.NewError(errors.ErrorUnauthorized, "already logged in on another device")
// 	// // 	log.Error().Msgf(err.Error())
// 	// // 	return nil, err
// 	// // }

// 	// b, _ := json.Marshal(ac)
// 	// _ = json.Unmarshal(b, &result)

// 	// pf, err := s.repo.GetProfileByID(ctx, ac.ID)
// 	// if err != nil {
// 	// 	log.Error().Msgf(err.Error())
// 	// 	return nil, err
// 	// }
// 	// b, _ = json.Marshal(pf)
// 	// _ = json.Unmarshal(b, &result)

// 	return &service.CheckAccessToken{
// 		AccountID:   userClaims.AccountInfo.ID,
// 		AppID:       userClaims.AccountInfo.AppID,
// 		AccountType: identity.AccountType(userClaims.AccountInfo.AccountType),
// 		RegisMode:   userClaims.AccountInfo.RegisMode,
// 		Status:      identity.AccountStatus(userClaims.AccountInfo.Status),
// 		EvStatus:    userClaims.AccountInfo.EvStatus,
// 		PvStatus:    userClaims.AccountInfo.PvStatus,
// 		Permission:  userClaims.AccountInfo.Permission,
// 		Password:    "****",
// 		Email:       userClaims.AccountInfo.Email,
// 		Phone:       userClaims.AccountInfo.Phone,
// 		Email_noti:  userClaims.AccountProfile.Email_noti,
// 		Phone_noti:  userClaims.AccountProfile.Phone_noti,
// 	}, nil
// }

// // ChangeEmail implements service.IdentityService.
// func (s *svc) ChangeEmail(ctx context.Context, id string, email string) (err error) {

// 	ac, _, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	// check email exist
// 	_, rnf, err := s.repo.GetByEmailForNormal(ctx, email)
// 	if err != nil {
// 		if !rnf {
// 			log.Error().Msgf("check exist by email failed: %v", err.Error())
// 			return err
// 		} else {
// 			err = nil
// 		}
// 	} else {
// 		err = errors.NewError(errors.ErrorConflict, "user email already exist.")
// 		log.Error().Msgf("user email already exist: %v ", email)
// 		return err
// 	}

// 	ac.Email = email
// 	ac.EvStatus = false

// 	// update
// 	err = s.repo.UpdateAccount(ctx, ac)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	userProfile, err := s.repo.GetProfileByID(ctx, ac.ID)
// 	if err != nil {
// 		return err
// 	}

// 	//寄驗證信
// 	err = s.sendEmailVerificationLinkAndCode(ctx, ac, userProfile, email, "email_verification")
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	// noti gorse
// 	gorseUserPatch := &gorse.UserPatch{}

// 	_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
// 	if err != nil {
// 		log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
// 		return err
// 	}

// 	return
// }

// // ChangePhone implements service.IdentityService.
// func (s *svc) ChangePhone(ctx context.Context, id string, phone string) (err error) {

// 	ac, _, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	// check phone

// 	if ac.Phone == phone {
// 		err = errors.NewError(errors.ErrorConflict, "user enter the same phone number")
// 		return err
// 	} else {
// 		ac.Phone = phone
// 		ac.PvStatus = false
// 	}

// 	// update
// 	err = s.repo.UpdateAccount(ctx, ac)
// 	if err != nil {
// 		log.Error().Msgf(err.Error())
// 		return err
// 	}

// 	//寄驗證信
// 	// to do

// 	// noti gorse
// 	gorseUserPatch := &gorse.UserPatch{}

//		_, err = s.gorseClient.UpdateUser(ctx, ac.ID, gorseUserPatch)
//		if err != nil {
//			log.Error().Msgf("gorse UpdateUser failed: %v", err.Error())
//			return err
//		}
//		return
//	}
func (s *svc) GenToken(ctx context.Context, accountInfo *repository.AccountInfo) (accessToken string, refreshToken string, err error) {
	tokenListKey := fmt.Sprintf("user-tokens:{%s}", accountInfo.ID)
	refreshTokenListKey := fmt.Sprintf("user-refresh-tokens:{%s}", accountInfo.ID)

	userProfile, err := s.repo.GetProfileByID(ctx, accountInfo.ID)
	if err != nil {
		return "", "", err
	}

	uuidGenerator := uuid.NewUUIDGenerator()

	accessToken, err = s.generateAndStoreToken(ctx, accountInfo, userProfile, "token", "tk", uuidGenerator, tokenListKey, s.accessTokenConfig.AccessTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateAndStoreToken(ctx, accountInfo, userProfile, "refresh_token", "rftk", uuidGenerator, refreshTokenListKey, s.accessTokenConfig.AccessTokenExpiresIn)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *svc) generateAndStoreToken(ctx context.Context, accountInfo *repository.AccountInfo, userProfile *repository.AccountProfile, tokenType string, prefix string, uuidGenerator uid.UIDGenerator, listKey string, expiresIn int) (string, error) {

	// userClaims := &service.UserClaims{
	// 	TokenType:      tokenType,
	// 	AccountInfo:    accountInfo,
	// 	AccountProfile: userProfile,
	// }

	// b, err := json.Marshal(userClaims)
	// if err != nil {
	// 	return "", errors.ErrorInternalError()
	// }

	tokenUUID := uuidGenerator.GenUID()
	// tokenSlotKey := fmt.Sprintf("%s-{%s}:%s", prefix, accountInfo.ID, tokenUUID)
	token := fmt.Sprintf("%s-%s:%s", prefix, accountInfo.ID, tokenUUID)

	// err = s.redisClusterClient.Set(ctx, tokenSlotKey, string(b), time.Minute*time.Duration(expiresIn)).Err()
	// if err != nil {
	// 	return "", errors.NewError(errors.ErrorInternalError, err.Error())
	// }

	// s.redisClusterClient.LPush(ctx, listKey, tokenSlotKey)
	// s.cleanupTokenList(ctx, tokenType, listKey, s.accessTokenConfig.MaxTokens)

	return token, nil
}

// func (s *svc) cleanupTokenList(ctx context.Context, tokenType string, listKey string, maxTokens int64) {
// 	l := log.Ctx(ctx)
// 	if maxTokens <= 0 {
// 		maxTokens = 100
// 	}

// 	deleteTokens, _ := s.redisClusterClient.LRange(ctx, listKey, maxTokens, -1).Result()
// 	if len(deleteTokens) > 0 {
// 		l.Info().Msgf("remove redundant %s tokens: %v", tokenType, deleteTokens)
// 		if err := s.redisClusterClient.Del(ctx, deleteTokens...).Err(); err != nil {
// 			l.Warn().Msgf("delete redundant tokens failed: %v", err.Error())
// 		}
// 	}

// 	s.redisClusterClient.LTrim(ctx, listKey, 0, maxTokens-1)
// 	s.redisClusterClient.Expire(ctx, listKey, time.Hour*24*7) // 列表保留 7 天
// }

// // 解析 token 獲取 claim data
// func (s *svc) GetTokenInfo(token string) (*service.JWTAuthClaims, error) {
// 	if token == "" {
// 		return nil, fmt.Errorf("no jwt token")
// 	}

// 	claim, err := s.parseToken(token)
// 	if err != nil {
// 		return nil, fmt.Errorf("bad jwt: %s", err)
// 	}

// 	return claim, nil
// }

// // 解析 jwt token
// func (s *svc) parseToken(token string) (*service.JWTAuthClaims, error) {
// 	jwtToken, err := jwt.ParseWithClaims(token, &service.JWTAuthClaims{}, func(token *jwt.Token) (i interface{}, e error) {
// 		return []byte(s.accessTokenConfig.JWTTokenKey), nil
// 	})
// 	if err == nil && jwtToken != nil {
// 		if claim, ok := jwtToken.Claims.(*service.JWTAuthClaims); ok && jwtToken.Valid {
// 			return claim, nil
// 		}
// 	}
// 	return nil, err
// }

// // 產生 jwt
// func (s *svc) GenerateJWT(id string, accountType int32, expired int64) (string, error) {
// 	claims := service.JWTAuthClaims{
// 		Fields: map[string]any{
// 			"ID":        id,
// 			"Type":      accountType,
// 			"TokenType": "user",
// 		},
// 		StandardClaims: jwt.StandardClaims{
// 			// Issuer:    Issuer,
// 			ExpiresAt: expired},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(s.accessTokenConfig.JWTTokenKey))
// 	if err != nil {
// 		return "", errors.NewError(errors.ErrorInternalError, "Sign token failed")
// 	}
// 	return tokenString, nil
// }

// func (s *svc) NormalizeTokenSlotKey(tokenKey string) (string, error) {
// 	parts := strings.SplitN(tokenKey, ":", 2)
// 	if len(parts) != 2 {
// 		return "", errors.NewErrorf(errors.ErrorInternalError, "invalid token key format: %s", tokenKey)
// 	}

// 	prefixParts := strings.SplitN(parts[0], "-", 2)
// 	if len(prefixParts) != 2 {
// 		return "", errors.NewErrorf(errors.ErrorInvalidInput, "invalid prefix-userID format: %s", parts[0])
// 	}

// 	normalizedKey := fmt.Sprintf("%s-{%s}:%s", prefixParts[0], prefixParts[1], parts[1])

// 	return normalizedKey, nil
// }
