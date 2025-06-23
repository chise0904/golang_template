package middleware

import (
	"github.com/chise0904/golang_template/pkg/trace"
	"github.com/chise0904/golang_template/pkg/uid"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var guid = uid.NewUIDGenerator(uid.GeneratorEnumUUID)

func NewRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = guid.GenUID()
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestID)

			logger := log.With().Str("request_id", requestID).Logger()
			ctx := logger.WithContext(c.Request().Context())
			ctx = trace.ContextWithXRequestID(ctx, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			// Set X-Request-Id header
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}
