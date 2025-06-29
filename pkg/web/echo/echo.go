package web

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/pkg/web"

	"github.com/chise0904/golang_template/pkg/web/echo/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEcho(config *web.Config, l fx.Lifecycle) *echo.Echo {
	echo.NotFoundHandler = EchoNotFoundHandler
	echo.MethodNotAllowedHandler = EchoNotAllowHandler

	e := echo.New()
	e.Validator = NewEchoValidator()

	if config.Mode == "release" {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	} else {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
	}

	e.HTTPErrorHandler = errors.EchoErrorHandler
	e.Pre(
		middleware.NewRequestIDMiddleware(),
	)
	e.Use(middleware.NewAccessLogMiddleware(config.RequestDump, config.MaxLogBodySize))
	e.Use(middleware.NewErrorHandlingMiddleware())
	e.Use(middleware.NewCORS())

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go e.StartServer(&http.Server{
				Addr: fmt.Sprintf(":%s", config.Port),
			})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			e.Shutdown(ctx)
			return nil
		},
	})
	return e
}

func EchoNotFoundHandler(c echo.Context) error {
	errors.EchoErrorHandler(errors.NewError(errors.ErrorPageNotFound, "page not found"), c)
	return nil
}

func EchoNotAllowHandler(c echo.Context) error {
	errors.EchoErrorHandler(errors.NewError(errors.ErrorNotAllow, "forbidden"), c)
	return nil
}

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return "j"
		}
		return name
	})
	return &EchoValidator{validator: v}
}

func (e *EchoValidator) Validate(i interface{}) error {
	return e.validator.Struct(i)
}
