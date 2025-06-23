package delivery

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// 	"gitlab.com/hsf-cloud/e-commerce/user-mgmt-server/service"
// 	"github.com/chise0904/golang_template/pkg/errors"
// )

// type createTokenRequest struct {
// 	GrantType    string `json:"grant_type" form:"grant_type" validate:"required"`       // password, authorization_code, refresh_token
// 	ClientID     string `json:"client_id" form:"client_id"  validate:"required"`        // application id
// 	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"` // application secret
// 	Email        string `json:"email" form:"email"`                                     // 用戶 email
// 	Phone        string `json:"phone" form:"phone"`                                     // 用戶 phone number
// 	Code         string `json:"code" form:"code"`                                       // ev_code pv_code
// 	Password     string `json:"password" form:"password"`                               // 用戶 密碼
// }

// type createTokenResponse struct {
// 	AccountID       string `json:"account_id"`
// 	Status          string `json:"status"`
// 	AccessToken     string `json:"access_token"`
// 	RefreshToken    string `json:"refresh_token"`
// 	TokenType       string `json:"token_type"`
// 	TokenExpireIN   int64  `json:"token_expire_in"`
// 	RefreshExpireIN int64  `json:"refresh_expire_in"`
// }

// func (h *handler) createToken(c echo.Context) error {

// 	connection := c.QueryParam("connection")
// 	if connection != "email" && connection != "sms" {
// 		return errors.NewError(errors.ErrorInvalidInput, "Invaild connection type")
// 	}

// 	req := &createTokenRequest{}
// 	err := c.Bind(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	err = c.Validate(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	in := &service.Login{
// 		Connection: connection,
// 		Email:      req.Email,
// 		Phone:      req.Phone,
// 		Code:       req.Code,
// 		Password:   req.Password,
// 	}
// 	resp := &createTokenResponse{}
// 	switch req.GrantType {
// 	case "password":
// 		id, status, acTk, rfTk, tkExp, rfExp, err := h.svc.CreateAccessTokenByPassword(c.Request().Context(), in)
// 		if err != nil {
// 			return err
// 		}
// 		resp = &createTokenResponse{
// 			AccountID:       id,
// 			Status:          status,
// 			AccessToken:     acTk,
// 			RefreshToken:    rfTk,
// 			TokenType:       "Bearer",
// 			TokenExpireIN:   tkExp,
// 			RefreshExpireIN: rfExp,
// 		}

// 	case "authorization_code":
// 		id, status, acTk, rfTk, tkExp, rfExp, err := h.svc.CreateAccessTokenByVeriCode(c.Request().Context(), in)
// 		if err != nil {
// 			return err
// 		}
// 		resp = &createTokenResponse{
// 			AccountID:       id,
// 			Status:          status,
// 			AccessToken:     acTk,
// 			RefreshToken:    rfTk,
// 			TokenType:       "Bearer",
// 			TokenExpireIN:   tkExp,
// 			RefreshExpireIN: rfExp,
// 		}
// 	case "refresh_token":
// 	default:
// 		return errors.NewError(errors.ErrorInvalidInput, "Invaild grant type")
// 	}
// 	return c.JSON(http.StatusCreated, resp)
// }

// type sendVerificationCodeRequest struct {
// 	Action string `json:"action" form:"action" validate:"required"` // login,setpassword
// 	Email  string `json:"email" form:"email"`                       // 用戶 email
// 	Phone  string `json:"phone" form:"phone"`                       // 用戶 phone number
// }

// func (h *handler) sendVerificationCode(c echo.Context) error {
// 	req := &sendVerificationCodeRequest{}
// 	err := c.Bind(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	err = c.Validate(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}

// 	connection := c.QueryParam("connection")
// 	switch connection {
// 	case "email":
// 		switch req.Action {
// 		case "login", "setpassword", "auth":
// 			in := &service.SendVerificationCode{
// 				Action: req.Action,
// 				Email:  req.Email,
// 			}
// 			err = h.svc.SendVerificationToEmail(c.Request().Context(), in)
// 			if err != nil {
// 				return err
// 			}
// 		default:
// 			return errors.NewError(errors.ErrorInvalidInput, "Invaild action type")
// 		}
// 	case "sms":
// 		switch req.Action {
// 		case "login", "setpassword", "auth":
// 			in := &service.SendVerificationCode{
// 				Action: req.Action,
// 				Phone:  req.Phone,
// 			}
// 			err = h.svc.SendVerificationToPhone(c.Request().Context(), in)
// 			if err != nil {
// 				return err
// 			}
// 		default:
// 			return errors.NewError(errors.ErrorInvalidInput, "Invaild action type")
// 		}
// 	default:
// 		return errors.NewError(errors.ErrorInvalidInput, "Invaild connection type")

// 	}
// 	return c.JSON(http.StatusCreated, "OK")
// }

// func (h *handler) checkAccessToken(c echo.Context) error {
// 	token := c.Param("accessToken")
// 	r, err := h.svc.CheckAccessToken(c.Request().Context(), token)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, r)
// }

// func (h *handler) googleLogin(c echo.Context) error {

// 	url := h.svc.CreateGoogleOAuthURL()

// 	// return c.JSON(http.StatusOK, url)

// 	return c.Redirect(http.StatusSeeOther, url)
// }

// func (h *handler) googleCallBack(c echo.Context) error {
// 	state := c.QueryParam("state")
// 	if state != service.GOOGLE_OAUTH_STATE {
// 		return errors.NewError(errors.ErrorUnauthorized, "Invaild state content")
// 	}

// 	code := c.QueryParam("code")
// 	// log.Debug().Msgf("code: %v", code)
// 	email, name, err := h.svc.GetUserInfoFromGoogle(c.Request().Context(), code)
// 	if err != nil {
// 		return err
// 	}

// 	resp := &createTokenResponse{}

// 	id, status, acTk, rfTk, tkExp, rfExp, err := h.svc.GoogleUserRegisAndLogin(c.Request().Context(), email, name)
// 	if err != nil {
// 		return err
// 	}
// 	resp = &createTokenResponse{
// 		AccountID:       id,
// 		Status:          status,
// 		AccessToken:     acTk,
// 		RefreshToken:    rfTk,
// 		TokenType:       "Bearer",
// 		TokenExpireIN:   tkExp,
// 		RefreshExpireIN: rfExp,
// 	}

// 	return c.JSON(http.StatusCreated, resp)
// }
