package middleware

import (
	// "github.com/chise0904/golang_template/pkg/trace"
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

			// 在任何 service 層或 handler 裡，只要有 ctx，都可以：
			// logger := log.Ctx(ctx)
			// logger.Info().Msg("something happened")
			// 這樣 log 裡面就會自動帶上：
			// {"level":"info","request_id":"abc-123","message":"something happened"}
			logger := log.With().Str("request_id", requestID).Logger()
			ctx := logger.WithContext(c.Request().Context())
			// ctx = trace.ContextWithXRequestID(ctx, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			// Set X-Request-Id header
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}
