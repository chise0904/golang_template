package middleware

import (
	"fmt"
	"runtime"

	"github.com/chise0904/golang_template/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// NewErrorHandlingMiddleware handles panic error
func NewErrorHandlingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stackError := make([]byte, 4096)
					runtime.Stack(stackError, true)
					requestID := c.Request().Header.Get(echo.HeaderXRequestID)
					customFields := map[string]interface{}{
						"url":         c.Request().RequestURI,
						"stack_error": string(stackError),
						"request_id":  requestID,
					}
					err, ok := r.(error)
					if !ok {
						if err == nil {
							err = fmt.Errorf("%v", r)
						} else {
							err = fmt.Errorf("%v", err)
						}
					}
					logger := log.With().Fields(customFields).Logger()
					logger.Error().Msgf("http: unknown error: %v", err)

					_ = c.JSON(500, errors.GetHttpError(errors.NewError(errors.ErrorInternalError, err.Error())))
				}
			}()
			return next(c)
		}
	}
}
