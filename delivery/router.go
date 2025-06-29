package delivery

import (
	"net/http"

	"github.com/chise0904/golang_template/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	svc service.IdentityService
}

func SetIdentityDelivery(e *echo.Echo, svc service.IdentityService) {

	h := &handler{
		svc: svc,
	}
	setRouter(e, h)
	setAuthRouter(e, h)
	setServerRouter(e, h)

}

func SetIdentityWebDelivery(e *echo.Echo) {

	setWebRouter(e)
}

func setRouter(e *echo.Echo, h *handler) {

	apiV1PublicAccountsGroup := e.Group("/public/apis/v1/identity/accounts")
	// apiV1PublicAccountsGroupRequireToken := e.Group("/public/apis/v1/identity/accounts")

	apiV1PublicAccountsGroup.Add("POST", "/register", h.registerAccount)
	apiV1PublicAccountsGroup.Add("POST", "/create", h.createAccount)
	// apiV1PublicAccountsGroup.Add("GET", "/emailVerification", h.emailVerification)

	// apiV1PublicAccountsGroup.Add("POST", "/passwords", h.setPassword)
	// //
	// apiV1PublicAccountsGroupRequireToken.Use(AccessTokenMiddleWare)

	// apiV1PublicAccountsGroupRequireToken.Add("DELETE", "/:accountID", h.deleteOne)
	// apiV1PublicAccountsGroupRequireToken.Add("GET", "/:accountID", h.getOne)
	// apiV1PublicAccountsGroupRequireToken.Add("GET", "", h.getAll)

	// apiV1PublicAccountsGroupRequireToken.Add("GET", "/profiles/:accountID", h.getProfile)
	// apiV1PublicAccountsGroupRequireToken.Add("PUT", "/profiles/:accountID", h.updateProfile)
	//
	//apiV1InternalAccountsGroup := e.Group("/internal/apis/v1/accounts")
}

func setAuthRouter(e *echo.Echo, h *handler) {
	// apiV1PublicOauthGroup := e.Group("/public/apis/v1/identity/oauth2")
	// apiV1InternalOauthGroup := e.Group("/internal/apis/v1/identity/oauth2")

	// apiV1PublicOauthGroup.Add("POST", "/token", h.createToken)                      // login
	// apiV1PublicOauthGroup.Add("POST", "/verificationCodes", h.sendVerificationCode) // for login,setpassword

	// apiV1PublicOauthGroup.Add("GET", "/google/login", h.googleLogin)
	// apiV1PublicOauthGroup.Add("GET", "/google/callback", h.googleCallBack)

	// apiV1InternalOauthGroup.Add("GET", "/token/:accessToken", h.checkAccessToken)
}

func setServerRouter(e *echo.Echo, h *handler) {
	apiV1PublicServerGroup := e.Group("/public/apis/v1/identity/server")
	apiV1PublicServerGroup.Add("GET", "/version", h.getVersion)

}

func setWebRouter(e *echo.Echo) {
	apiV1PublicWebGroup := e.Group("/public/apis/v1/identity-web")
	apiV1PublicWebGroup.Static("/emailVerification/en", "../template/html/en")
	apiV1PublicWebGroup.Static("/emailVerification/tw", "../template/html/tw")

	apiV1PublicWebGroup.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.Start(":8081")

}

// func AccessTokenMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		token := c.Request().Header.Get("access_token")
// 		if token == "" {
// 			token = c.QueryParam("access_token")
// 		}

// 		_, _, err := access_token.GetExpired(token)
// 		if err != nil {
// 			return err
// 		}
// 		return next(c)
// 	}
// }
