package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chise0904/golang_template/pkg/errors"
	"gorm.io/plugin/soft_delete"
)

type GetAccountListOpts struct {
	IDs          []string
	AccountTypes []int32
	Sort         string // ASC.DESC
	By           string // id.email.time
	Page         uint32 // offset
	Perpage      uint32 // limit, default 25
	Filter       string
}
type AccountInfo struct {
	ID          string `json:"account_id" gorm:"id"`
	AppID       string `json:"app_id" gorm:"app_id"`
	AccountType int32  `json:"account_type" gorm:"account_type"` // admin,user,vendor...etc
	RegisMode   string `json:"regis_mode" gorm:"regis_mode"`     // default,google...etc
	Status      int32  `json:"status" gorm:"status"`             //

	EvStatus   bool       `json:"ev_status" gorm:"ev_status"`   // email 驗證狀態
	PvStatus   bool       `json:"pv_status" gorm:"pv_status"`   // phone 驗證狀態
	Permission Permission `json:"permission" gorm:"permission"` // 權限

	Password string `json:"password" gorm:"password"` // user password
	Email    string `json:"email" gorm:"email"`       // user email
	Phone    string `json:"phone" gorm:"phone"`       // user phone (optional)

	LoginAt   int64                 `json:"login_at" gorm:"login_at"`               // timestamp ms
	LogoutAt  int64                 `json:"logout_at" gorm:"logout_at"`             // timestamp ms
	CreatedAt int64                 `json:"created_at" gorm:"autoCreateTime:false"` // timestamp ms
	UpdatedAt int64                 `json:"updated_at" gorm:"autoUpdateTime:false"` // timestamp ms
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"softDelete:milli"`
}

type SimpAccountInfo struct {
	ID          string `json:"account_id" gorm:"id"`
	AccountType string `json:"account_type" gorm:"account_type"`       // admin,user...etc
	RegisMode   string `json:"regis_mode" gorm:"regis_mode"`           // default,google...etc
	Status      string `json:"status" gorm:"status"`                   //
	Email       string `json:"email" gorm:"email"`                     // user email
	CreatedAt   int64  `json:"created_at" gorm:"autoCreateTime:false"` // timestamp ms
}

type AccountProfile struct {
	AccountID       string       `json:"account_id" gorm:"account_id"`
	UserName        string       `json:"user_name" gorm:"user_name"`
	Icon            []byte       `json:"icon" gorm:"icon"`
	Description     string       `json:"description" gorm:"description"`
	Gender          string       `json:"gender" gorm:"gender"`
	Birthday        Date         `json:"birthday" gorm:"birthday"`
	Job             string       `json:"job" gorm:"job"`
	Country         string       `json:"country" gorm:"country"`
	City            string       `json:"city" gorm:"city"`
	District        string       `json:"district" gorm:"district"`
	ZipCode         string       `json:"zip_code" gorm:"zip_code"`
	Address         string       `json:"address" gorm:"address"`
	ShippingAddress AddressArray `json:"shipping_address" gorm:"shipping_address"`

	Language   string                `json:"language" gorm:"language"`
	Phone_noti bool                  `json:"phone_noti" gorm:"phone_noti"`
	Email_noti bool                  `json:"email_noti" gorm:"email_noti"`
	CreatedAt  int64                 `json:"created_at" gorm:"autoCreateTime:false"` // timestamp ms
	UpdatedAt  int64                 `json:"updated_at" gorm:"autoUpdateTime:false"` // timestamp ms
	DeletedAt  soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"softDelete:milli"`
}

type DBProfileOptions struct {
	AccountID       string       `json:"account_id" gorm:"account_id"`
	UserName        string       `json:"user_name" gorm:"user_name"`
	Icon            []byte       `json:"icon" gorm:"icon"`
	Description     string       `json:"description" gorm:"description"`
	Gender          string       `json:"gender" gorm:"gender"`
	Birthday        *Date        `json:"birthday" gorm:"birthday"`
	Job             string       `json:"job" gorm:"job"`
	Country         string       `json:"country" gorm:"country"`
	City            string       `json:"city" gorm:"city"`
	District        string       `json:"district" gorm:"district"`
	ZipCode         string       `json:"zip_code" gorm:"zip_code"`
	Address         string       `json:"address" gorm:"address"`
	ShippingAddress AddressArray `json:"shipping_address" gorm:"shipping_address"`
	Language        string       `json:"language" gorm:"language"`
	Phone_noti      *bool        `json:"phone_noti" gorm:"phone_noti"`
	Email_noti      *bool        `json:"email_noti" gorm:"email_noti"`
}

type Date struct {
	Day   int32 `json:"day"`
	Month int32 `json:"month"`
	Year  int32 `json:"year"`
}

type AddressArray []Address
type Address struct {
	Type     string `json:"type"`
	Country  string `json:"country"`
	City     string `json:"city"`
	District string `json:"district"`
	ZipCode  string `json:"zip_code"`
	Address  string `json:"address"`
	StoreID  string `json:"store_id"`
}

type Permission struct {
	CanAccessCrossAccount bool `json:"can_access_cross_account"`
	ProductRead           bool `json:"product_read"`
	ProductRewrite        bool `json:"product_rewrite"`
	OrderRead             bool `json:"order_read"`
	OrderRewrite          bool `json:"order_rewrite"`
	SubscribeEmail        bool `json:"subscribe_email"`
	CoMarketing           bool `json:"co_marketing"`
}

type VerificationCode struct {
	AccountID string `json:"account_id" gorm:"account_id"`
	Action    string `json:"action" gorm:"action"`
	Code      string `json:"code" gorm:"code"`
	Token     string `json:"token" gorm:"token"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:false"` // timestamp ms
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime:false"` // timestamp ms
	DeletedAt int64  `json:"deleted_at,omitempty" gorm:"softDelete:milli"`
}

func (s *AccountInfo) GetTableName() string {
	return "account_info"
}

func (s *AccountProfile) GetTableName() string {
	return "user_profile"
}

func (s *VerificationCode) GetTableName() string {
	return "verification_code"
}

func (a *AddressArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.NewErrorf(errors.ErrorInternalError, "Failed to unmarshal JSONB value=%v", value)
	}
	err := json.Unmarshal(bytes, a)
	return err
}

func (a AddressArray) Value() (driver.Value, error) {
	json.Marshal(a)
	return json.Marshal(a)
}
func (s *Permission) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.NewErrorf(errors.ErrorInternalError, "Failed to unmarshal JSONB value=%v", value)
	}
	err := json.Unmarshal(bytes, s)
	return err
}

func (p Permission) Value() (driver.Value, error) {
	json.Marshal(p)
	return json.Marshal(p)
}

func (d *Date) Scan(value interface{}) error {
	switch s := value.(type) {
	case time.Time:
		*d = Date{Year: int32(s.Year()), Month: int32(s.Month()), Day: int32(s.Day())}
	case []byte:
		t, err := time.Parse("2006-01-02", string(s))
		if err != nil {
			return errors.NewError(errors.ErrorInternalError, err.Error())
		}
		*d = Date{Year: int32(t.Year()), Month: int32(t.Month()), Day: int32(t.Day())}
	case string:
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return errors.NewError(errors.ErrorInternalError, err.Error())
		}
		*d = Date{Year: int32(t.Year()), Month: int32(t.Month()), Day: int32(t.Day())}
	default:
		return errors.NewErrorf(errors.ErrorInternalError, "date: Unsupport scanning type %T", value)
	}
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return fmt.Sprintf("%4d-%02d-%02d", d.Year, d.Month, d.Day), nil
}

type ListProfilesFilter struct {
	AccountIDs []string
	Page       uint32 // offset
	Perpage    uint32 // limit
}
