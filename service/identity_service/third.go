package impl

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/pkg/uid/xid"
	"github.com/chise0904/golang_template/proto/pkg/identity"
	"github.com/chise0904/golang_template/repository"
	"github.com/chise0904/golang_template/service"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// CreateGoogleOAuthURL implements service.IdentityService
func (s *svc) CreateGoogleOAuthURL() string {
	config := &oauth2.Config{

		ClientID:     s.oauthConfig.GoogleClientID,
		ClientSecret: s.oauthConfig.GoogleSecretKey,
		RedirectURL:  s.oauthConfig.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// 這行會產生一個 完整的 Google 授權 URL，裡面會包含：
	// client_id
	// redirect_uri
	// scope
	// response_type=code
	// state ← 很重要！
	//
	// state 是一個隨機字串（或你的 session id）
	// 你傳給 Google，它會原封不動回傳
	// 你就能驗證「這是我發出去的請求，不是 CSRF 攻擊」
	// 如果 state 不一致，就拒絕這次登入流程
	//
	// 一個產生出來的 URL 長這樣：
	// https://accounts.google.com/o/oauth2/auth?
	// client_id=your-client-id
	// &redirect_uri=https://your-app.com/auth/callback
	// &scope=email+profile
	// &response_type=code
	// &state=random123
	// 使用者點這個 URL，Google 顯示授權畫面，授權成功後會 redirect 到：
	// https://your-app.com/auth/callback?code=abc123&state=random123
	// 你就可以用這個 code 換 access token，實作登入功能。
	return config.AuthCodeURL(s.oauthConfig.GoogleOAuthState) // 防偽用 state param, google 會原封不動回傳, 可供驗證
}

// GetUserInfoFromGoogle implements service.IdentityService
func (s *svc) GetUserInfoFromGoogle(ctx context.Context, code string) (email string, name string, err error) {
	// var google_config *oauth2.Config

	google_config := &oauth2.Config{

		ClientID:     s.oauthConfig.GoogleClientID,
		ClientSecret: s.oauthConfig.GoogleSecretKey,
		RedirectURL:  s.oauthConfig.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// 這是標準 OAuth 流程的第二步：
	// code 是 Google redirect 回來時附帶的授權碼
	// 你用它換取 access_token（t 代表 Token）
	t, err := google_config.Exchange(ctx, code)
	if err != nil {
		return email, name, errors.NewError(errors.ErrorUnauthorized, err.Error())
	}
	client := google_config.Client(context.TODO(), t)
	// 這會取得 JSON 格式的使用者資料（email、name、picture 等）
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return email, name, errors.NewError(errors.ErrorInvalidInput, err.Error())
	}
	defer userInfo.Body.Close()

	info, err := io.ReadAll(userInfo.Body)
	if err != nil {
		return email, name, errors.NewError(errors.ErrorInternalError, err.Error())
	}
	var user service.GoogleUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		return email, name, errors.NewError(errors.ErrorInternalError, err.Error())
	}

	return user.Email, user.Name, nil
}

// GoogleUserRegisAndLogin implements service.IdentityService
func (s *svc) GoogleUserRegisAndLogin(ctx context.Context, email string, name string) (id string, status int32, token, refresh_token string, token_expired, refresh_expired int64, err error) {
	token_expired = time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.AccessTokenExpiresIn)).UnixMilli()
	refresh_expired = time.Now().Add(time.Minute * time.Duration(s.accessTokenConfig.RefreshTokenExpiresIn)).UnixMilli()
	resp := &repository.AccountInfo{}
	nowMS := time.Now().UnixMilli()
	ac, rnf, err := s.repo.GetByEmailForThirdParty(ctx, email, "google")
	if err != nil {
		if !rnf {
			log.Error().Msgf("check exist by email failed: %v", err.Error())
			return "", 0, "", "", token_expired, refresh_expired, err
		} else { // 查無帳號 -> 則註冊
			acModel := &repository.AccountInfo{
				AccountType: int32(identity.AccountType_AccountType_USER),
				Status:      int32(identity.AccountStatus_AccountStatus_ENABLED),
				EvStatus:    true,
				Email:       email,
			}
			// 生成ID
			acXID := xid.NewXIDGenerator().GenUID()
			acModel.ID = acXID
			acModel.RegisMode = "GOOGLE"
			// 創建新帳號
			err = s.repo.UpsertAccount(ctx, acModel)
			if err != nil {
				log.Error().Msgf(err.Error())
				return "", 0, "", "", token_expired, refresh_expired, err
			}
			// 創建profile
			acProfile := &repository.AccountProfile{
				AccountID: acXID,
				UserName:  name,
				CreatedAt: nowMS,
			}
			err = s.repo.UpsertProfile(ctx, acProfile)
			if err != nil {
				log.Error().Msgf(err.Error())
				return "", 0, "", "", token_expired, refresh_expired, err
			}
			ac = acModel
		}
	}
	// 如果帳號info已存在，則無需 create，直接login
	// 若本不存在則接續註冊行為之後 login
	if ac.Email == email && ac.Status == int32(identity.AccountStatus_AccountStatus_BLOCKED) {
		return "", 0, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorNotAllow, "account is blocked")
		// return id, status, accessToken, expired, errors.NewError(errors.ErrorConflict, "User's email has insufficient permissions.")
	} else if ac.Email == email && ac.Status == int32(identity.AccountStatus_AccountStatus_ENABLED) {

		// gen token & refresh token
		token, err = s.GenerateJWT(ac.ID, ac.AccountType, token_expired)
		if err != nil {
			log.Error().Msgf("gen token failed: %v", err.Error())
			return "", 0, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		}
		refresh_token, err = s.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		if err != nil {
			log.Error().Msgf("gen token failed: %v", err.Error())
			return "", 0, "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		}
		// isExp, err := access_token.IsExpired(ac.RefreshToken)
		// if isExp || err != nil {
		// 	refresh_expired := time.Now().Add(time.Minute * time.Duration(access_token.REFRESH_TOKEN_EXPIRES_IN)).UnixMilli()
		// 	ac.RefreshToken, err = access_token.GenerateJWT(ac.ID, ac.AccountType, refresh_expired)
		// 	if err != nil {
		// 		log.Error().Msgf("gen token failed: %v", err.Error())
		// 		return "", "", "", "", token_expired, refresh_expired, errors.NewError(errors.ErrorInternalError, err.Error())
		// 	}
		// }

		ac.LoginAt = nowMS
		ac.UpdatedAt = nowMS

		err = s.repo.UpdateAccount(ctx, ac)
		if err != nil {
			log.Error().Msgf("update account for login failed: %v", err.Error())
			return "", 0, "", "", token_expired, refresh_expired, err
		}
		resp = ac
	}

	return resp.ID, resp.Status, token, refresh_token, token_expired, refresh_expired, err
}
