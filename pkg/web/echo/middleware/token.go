package middleware

import (
	"github.com/chise0904/golang_template/pkg/auth"
	"github.com/labstack/echo/v4"
)

type JWTAuthHandler interface {
	TokenVerify(c echo.Context) (*auth.UserClaims, error)
}

func JWTAuthMiddlewareFunc(h JWTAuthHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		//check access token...
		return func(c echo.Context) error {
			claims, err := h.TokenVerify(c)
			if err != nil {
				return err
			}
			ctx := auth.UserClaimsWithContext(c.Request().Context(), claims)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
