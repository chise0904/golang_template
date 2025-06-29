package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
	// 	"gitlab.com/hsf-cloud/e-commerce/user-mgmt-server/repository"
	// 	"gitlab.com/hsf-cloud/e-commerce/user-mgmt-server/service"

	"github.com/chise0904/golang_template/pkg/errors"
	// "github.com/chise0904/golang_template/pkg/web"
	"github.com/chise0904/golang_template/proto/pkg/common"
	"github.com/chise0904/golang_template/proto/pkg/identity"
)

type registerAccountRequest struct {
	Email    string `json:"email" form:"email" validate:"required"`         // 用戶 email
	Password string `json:"password" form:"password" validate:"required"`   // 用戶 密碼
	UserName string `json:"user_name" form:"user_name" validate:"required"` // 用戶 名稱
}
type registerAccountResponse struct {
	Email string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	// Phone string `form:"phone,omitempty" json:"phone,omitempty" xml:"phone,omitempty"`
	// UserName string `form:"user_name,omitempty" json:"user_name,omitempty" xml:"user_name,omitempty"`
	Href string `form:"href" json:"href" xml:"href"`
}

func (h *handler) registerAccount(c echo.Context) error {
	req := &registerAccountRequest{}
	err := c.Bind(req)
	if err != nil {
		return errors.NewError(errors.ErrorInvalidInput, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return errors.NewError(errors.ErrorInvalidInput, err.Error())
	}

	r, err := h.svc.RegisterAccount(c.Request().Context(), req.Email, req.Password, req.UserName, identity.AccountType_AccountType_USER, nil)
	if err != nil {
		return err
	}
	out := &registerAccountResponse{
		Email: r.Email,
		Href:  "",
	}
	return c.JSON(http.StatusCreated, out)

}

type createAccountRequest struct {
	Email    string `json:"email" form:"email" validate:"required"`         // 用戶 email
	Password string `json:"password" form:"password" validate:"required"`   // 用戶 密碼
	UserName string `json:"user_name" form:"user_name" validate:"required"` // 用戶 名稱
}

type createAccountResponse struct {
	Email string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	// Phone string `form:"phone,omitempty" json:"phone,omitempty" xml:"phone,omitempty"`
	UserName string `form:"user_name,omitempty" json:"user_name,omitempty" xml:"user_name,omitempty"`
	Href     string `form:"href" json:"href" xml:"href"`
}

func (h *handler) createAccount(c echo.Context) error {
	req := &createAccountRequest{}
	err := c.Bind(req)
	if err != nil {
		return errors.NewError(errors.ErrorInvalidInput, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return errors.NewError(errors.ErrorInvalidInput, err.Error())
	}

	// 舊的HTTP創建帳號方式
	// r, err := h.svc.CreateAccount(c.Request().Context(), req.Email, req.Password, req.UserName, "user")

	// 創建 CreateAccountRequest
	createReq := &identity.CreateAccountRequest{
		Email:       req.Email,
		Password:    req.Password,
		UserName:    req.UserName,
		AccountType: identity.AccountType_AccountType_USER,
		RegisMode:   "NORMAL",
		Status:      identity.AccountStatus_AccountStatus_ENABLED,
		Permission: &identity.Permission{
			CanAccessCrossAccount:     false,
			CanReadProduct:            true,
			CanModifyProduct:          true,
			CanReadOrder:              true,
			CanModifyOrder:            true,
			CanReceiveEmails:          false,
			CanParticipateInMarketing: false,
		},
		EmailNoti: common.BoolType_False,
		PhoneNoti: common.BoolType_False,
	}

	r, err := h.svc.CreateAccount(c.Request().Context(), createReq)
	if err != nil {
		return err
	}
	out := &createAccountResponse{
		Email: r.Email,
		Href:  "",
	}
	return c.JSON(http.StatusCreated, out)
}

// type setPasswordRequest struct {
// 	OldPassword string `json:"old_password" form:"old_password"`             // 用戶 舊密碼 - by password
// 	Password    string `json:"password" form:"password" validate:"required"` // 用戶 密碼
// 	Phone       string `json:"phone" form:"phone"`                           // 用戶 phone - by verify code
// 	Email       string `json:"email" form:"email"`                           // 用戶 email - by verify code
// 	Code        string `json:"code" form:"code"`                             // verify code
// }

// func (h *handler) setPassword(c echo.Context) error {

// 	by := c.QueryParam("by")

// 	req := &setPasswordRequest{}
// 	err := c.Bind(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	err = c.Validate(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}

// 	switch by {
// 	case "password":
// 		token := c.Request().Header.Get("access_token")
// 		if token == "" {
// 			token = c.QueryParam("access_token")
// 		}

// 		in := &service.SetPassword{
// 			AccessToken: token,
// 			OldPassword: req.OldPassword,
// 			Password:    req.Password,
// 		}
// 		err = h.svc.SetPasswordByPassword(c.Request().Context(), by, in)
// 		if err != nil {
// 			return err
// 		}

// 	case "email", "phone":

// 		in := &service.SetPassword{
// 			Phone:    req.Phone,
// 			Email:    req.Email,
// 			Code:     req.Code,
// 			Password: req.Password,
// 		}
// 		err = h.svc.SetPasswordByVerifyCode(c.Request().Context(), by, in)
// 		if err != nil {
// 			return err
// 		}
// 	default:
// 		return errors.NewError(errors.ErrorInvalidInput, "Invaild set password type")
// 	}

// 	return c.JSON(http.StatusOK, "OK")
// }

// func (h *handler) emailVerification(c echo.Context) error {
// 	token := c.QueryParam("access_token")

// 	_, url, _ := h.svc.EmailVerification(c.Request().Context(), token)

// 	return c.Redirect(http.StatusFound, url)
// }

// func (h *handler) getOne(c echo.Context) error {

// 	id := c.Param("accountID")

// 	token := c.Request().Header.Get("access_token")
// 	if token == "" {
// 		token = c.QueryParam("access_token")
// 	}

// 	r, err := h.svc.GetOne(c.Request().Context(), token, id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, r)
// }

// type getAllRequest struct {
// 	By          string `query:"by"`
// 	Sort        string `query:"sort"`
// 	Page        uint32 `query:"page"`
// 	Perpage     uint32 `query:"perpage"`
// 	Filter      string `query:"filter"`
// 	AccessToken string `query:"access_token"`
// }

// type getAllResponse struct {
// 	Meta     web.ResponsePayLoadMetaData `json:"meta"`
// 	Accounts []repository.AccountInfo    `json:"accounts"`
// }

// type getAllSimpResponse struct {
// 	Meta     web.ResponsePayLoadMetaData  `json:"meta"`
// 	Accounts []repository.SimpAccountInfo `json:"accounts"`
// }

// func (h *handler) getAll(c echo.Context) error {

// 	req := &getAllRequest{}
// 	err := c.Bind(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}

// 	token := c.Request().Header.Get("access_token")
// 	if token == "" {
// 		token = req.AccessToken
// 	}

// 	opts := &repository.GetAccountListOpts{}

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

// 	admin_request, simp_data, detail_data, pg, err := h.svc.GetAll(c.Request().Context(), token, opts)
// 	if err != nil {
// 		return err
// 	}

// 	if !admin_request {
// 		r := &getAllSimpResponse{
// 			Meta: web.ResponsePayLoadMetaData{
// 				Pagination: pg,
// 			},
// 			Accounts: simp_data,
// 		}
// 		return c.JSON(http.StatusOK, r)
// 	} else {
// 		r := &getAllResponse{
// 			Meta: web.ResponsePayLoadMetaData{
// 				Pagination: pg,
// 			},
// 			Accounts: detail_data,
// 		}
// 		return c.JSON(http.StatusOK, r)
// 	}

// }

// func (h *handler) deleteOne(c echo.Context) error {
// 	id := c.Param("accountID")
// 	token := c.Request().Header.Get("access_token")
// 	if token == "" {
// 		token = c.QueryParam("access_token")
// 	}

// 	err := h.svc.DeleteAccount(c.Request().Context(), token, id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, "OK")
// }

// // type getProfileResponse struct {
// // 	ID        string `json:"id"`
// // 	UserName  string `json:"user_name"`
// // 	Icon      string `json:"icon"`
// // 	CreatedAt int64  `json:"created_at"` // timestamp ms
// // 	UpdatedAt int64  `json:"updated_at"` // timestamp ms
// // }

// func (h *handler) getProfile(c echo.Context) error {
// 	id := c.Param("accountID")
// 	token := c.Request().Header.Get("access_token")
// 	if token == "" {
// 		token = c.QueryParam("access_token")
// 	}

// 	r, err := h.svc.GetUserProfile(c.Request().Context(), token, id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, r)
// }

// type updateProfileRequest struct {
// 	Icon        []byte          `json:"icon"`
// 	UserName    string          `json:"user_name" validate:"required"`
// 	Description string          `json:"description" validate:"required"`
// 	Gender      string          `json:"gender" validate:"required"`
// 	Birthday    repository.Date `json:"birthday" validate:"required"`
// 	Job         string          `json:"job" validate:"required"`
// 	Country     string          `json:"country" validate:"required"`
// 	City        string          `json:"city" validate:"required"`
// 	Address     string          `json:"address" validate:"required"`
// }

// func (h *handler) updateProfile(c echo.Context) error {
// 	req := &updateProfileRequest{}
// 	err := c.Bind(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	err = c.Validate(req)
// 	if err != nil {
// 		return errors.NewError(errors.ErrorInvalidInput, err.Error())
// 	}
// 	id := c.Param("accountID")
// 	token := c.Request().Header.Get("access_token")
// 	if token == "" {
// 		token = c.QueryParam("access_token")
// 	}
// 	// b := fmt.Sprintf("%v-%v-%v", req.Birthday.Year, req.Birthday.Month, req.Birthday.Day)
// 	pf := &repository.AccountProfile{
// 		AccountID:   id,
// 		UserName:    req.UserName,
// 		Icon:        req.Icon,
// 		Description: req.Description,
// 		Gender:      req.Gender,
// 		Birthday: repository.Date{
// 			Day:   req.Birthday.Day,
// 			Month: req.Birthday.Month,
// 			Year:  req.Birthday.Year,
// 		},
// 		Job:     req.Job,
// 		Country: req.Country,
// 		City:    req.City,
// 		Address: req.Address,
// 	}
// 	err = h.svc.UpdateUserProfile(c.Request().Context(), token, pf)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, "OK")
// }
