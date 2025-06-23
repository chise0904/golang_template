package errors

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type ErrorView struct {
	Category Category `json:"category"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
}

func GetHttpError(err error) ErrorView {
	e := errors.Cause(err)
	if _err, ok := e.(*_error); ok {
		return ErrorView{
			Message:  err.Error(),
			Code:     _err.code,
			Category: _err.category,
		}
	}
	return ErrorView{
		Message: err.Error(),
		Code:    "",
	}
}

// EchoErrorHandler responds error response according to given error.
func EchoErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	echoErr, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoErr.Code, echoErr)
		return
	}

	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok || _err == nil {
		_ = c.JSON(http.StatusInternalServerError, ErrorInternalError)
		return
	}

	_ = c.JSON(int(_err.category), GetHttpError(_err))
}

func ConvertGrpcErrToHttpErr(err error) error {
	if err == nil {
		return nil
	}
	s := status.Convert(err)
	if s == nil {
		return NewError(ErrorInternalError, "")
	}
	interErr := _error{}

	jerr := json.Unmarshal([]byte(s.Message()), &interErr)
	if jerr != nil {
		return convertGrpcStatusToHttpErr(s)
	}
	interErr.grpcCode = s.Code()

	return WithStack(&interErr)
}

// ConvertErrorToGrpcErr Convert _error to grpc error
func ConvertErrorToGrpcErr(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return status.Error(ErrorInternalError().grpcCode, err.Error())
	}

	b, _ := json.Marshal(_err)
	return status.Error(_err.grpcCode, string(b))
}

func convertGrpcStatusToHttpErr(s *status.Status) error {
	var httpError error
	switch s.Code() {
	case Unknown:
		httpError = NewError(ErrorInternalError, s.Message())
	case InvalidArgument:
		httpError = NewError(ErrorInvalidInput, s.Message())
	case NotFound:
		httpError = NewError(ErrorResourceNotFound, s.Message())
	case AlreadyExists:
		httpError = NewError(ErrorConflict, s.Message())
	case PermissionDenied:
		httpError = NewError(ErrorNotAllow, s.Message())
	case Unauthenticated:
		httpError = NewError(ErrorUnauthorized, s.Message())
	case OutOfRange:
		httpError = NewError(ErrorInvalidInput, s.Message())
	case Internal:
		httpError = NewError(ErrorInternalError, s.Message())
	case DataLoss:
		httpError = NewError(ErrorInternalError, s.Message())

	default:
		httpError = NewError(ErrorInternalError, s.Message())
	}
	return httpError
}
