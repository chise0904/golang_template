package service

// import (
// 	 "time"
// )

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

type SendVerificationCode struct {
	Action string `json:"action" validate:"required"` // login, setpassword, auth
	Email  string `json:"email"`                      // 用戶 email
	Phone  string `json:"phone"`                      // 用戶 phone number
}

type Version struct {
	Version      string `json:"version"`
	BuildTime    string `json:"buildTime"`
	Sha1ver      string `json:"sha1ver"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

type Login struct {
	Connection string `json:"connection"` // email,sms
	Email      string `json:"email"`      // 用戶 email
	Phone      string `json:"phone"`      // 用戶 phone number
	Code       string `json:"code"`       // ev_code pv_code
	Password   string `json:"password"`   // 用戶 密碼
}

type SetPassword struct {
	AccessToken string `json:"access_token"`
	OldPassword string `json:"old_password"` // 用戶 舊密碼
	Phone       string `json:"phone"`        // 用戶 phone
	Email       string `json:"email"`        // 用戶 email
	Code        string `json:"code"`         // verify code
	Password    string `json:"password"`     // 用戶 密碼
}

type AccessTokenConfig struct {
	AccessTokenExpiresIn  int    `mapstructure:"access_token_expires_in"`
	RefreshTokenExpiresIn int    `mapstructure:"refresh_token_expires_in"`
	JWTTokenKey           string `mapstructure:"token_key"`
	MaxTokens             int64  `mapstructure:"max_tokens"`
	// Issuer                string `mapstructure:"issuer"`
}

type NotifyServiceConfig struct {
	Setting EmailSetting `mapstructure:"email_setting"`
	Content EmailContent `mapstructure:"email_content"`
}

type EmailSetting struct {
	TokenExpiresIn int    `mapstructure:"token_expires_in"`
	CodeExpiresIn  int    `mapstructure:"code_expires_in"`
	TokenKey       string `mapstructure:"token_key"`

	MailgunSenderMail    string `mapstructure:"mailgun_sender_mail"`
	MailgunSendingDomain string `mapstructure:"mailgun_sending_domain"`
	MailgunSendingApiKey string `mapstructure:"mailgun_sending_api_key"`
}
type EmailContent struct {
	CustomerServiceContactInformation string `mapstructure:"customer_service_contact_information"`
	VerificationURL                   string `mapstructure:"verification_url"`
	RegistryOrChangeTitle             string `mapstructure:"registry_or_change_title"`
	RegistryOrChangeTemplate          string `mapstructure:"registry_or_change_template"`
	FindPasswordTitle                 string `mapstructure:"find_password_title"`
	FindPasswordTemplate              string `mapstructure:"find_password_template"`
	LoginByCodeTitle                  string `mapstructure:"login_by_code_title"`
	LoginByCodeTemplate               string `mapstructure:"login_by_code_template"`
	FromName                          string `mapstructure:"from_name"`
	ServiceName                       string `mapstructure:"service_name"`
}

type Recipient struct {
	Name      string
	EmailAddr string
}
