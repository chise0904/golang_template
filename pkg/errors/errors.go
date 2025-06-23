package errors

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type _error struct {
	category Category
	code     string
	message  string
	indx     int
	top      *int
	grpcCode codes.Code
}

type withMessage struct {
	cause error
	err   error
}

var WithStack = errors.WithStack

func GetErrorCode(err error) string {
	if e, ok := errors.Cause(err).(*_error); ok {
		return e.code
	} else {
		return ""
	}
}

func (err *_error) Error() string {

	if err.indx != *err.top {
		return err.message
	}

	return fmt.Sprintf("[%s] %s", err.code, err.message)
}

func (w *withMessage) Error() string {

	err, ok := w.err.(*_error)
	if !ok {
		return fmt.Sprintf("%s: %v", w.err.Error(), w.cause)
	}
	if err.indx == *err.top {
		return fmt.Sprintf("[%s] %s: %v", err.code, err.message, w.cause)
	}
	return fmt.Sprintf("%s: %v", err.message, w.cause)

}

func (w *withMessage) Cause() error {
	if w.cause == nil {
		return nil
	}
	return w.cause
}

func NewError(err baseError, message string) error {
	_err := err()
	return WithStack(&_error{
		category: _err.category,
		code:     _err.code,
		grpcCode: _err.grpcCode,
		message:  message,
		top:      new(int),
	})

}

func NewErrorf(err baseError, formate string, a ...any) error {
	_err := err()
	return WithStack(&_error{
		category: _err.category,
		code:     _err.code,
		grpcCode: _err.grpcCode,
		message:  fmt.Sprintf(formate, a...),
		top:      new(int),
	})

}

func Wrap(err error, message string) error {
	switch e := errors.Cause(err).(type) {
	case *_error:
		if e.top == nil {
			e.top = new(int)
		}
		*e.top++
		w := &withMessage{
			cause: err,
			err: &_error{
				indx:     *e.top,
				top:      e.top,
				category: e.category,
				code:     e.code,
				grpcCode: e.grpcCode,
				message:  message,
			},
		}

		return w
	case *withMessage:

		if e, ok := e.err.(*_error); ok {
			if e.top == nil {
				e.top = new(int)
			}
			*e.top++
			return &withMessage{
				cause: err,
				err: &_error{
					indx:     *e.top,
					top:      e.top,
					category: e.category,
					code:     e.code,
					grpcCode: e.grpcCode,
					message:  message,
				},
			}
		}
		return &withMessage{
			cause: err,
			err:   errors.Errorf(message),
		}
	default:
		return &withMessage{
			cause: err,
			err:   errors.Errorf(message),
		}
	}

}

func Wrapf(err error, formate string, a ...any) error {

	switch e := errors.Cause(err).(type) {
	case *_error:
		if e.top == nil {
			e.top = new(int)
		}
		*e.top++
		w := WithStack(&withMessage{
			cause: err,
			err: &_error{
				indx:     *e.top,
				top:      e.top,
				category: e.category,
				code:     e.code,
				grpcCode: e.grpcCode,
				message:  fmt.Sprintf(formate, a...),
			},
		})
		return w
	case *withMessage:

		if _e, ok := e.err.(*_error); ok {
			*_e.top++
			return WithStack(&withMessage{
				cause: err,
				err: &_error{
					indx:     *_e.top,
					top:      _e.top,
					category: _e.category,
					code:     _e.code,
					grpcCode: _e.grpcCode,
					message:  fmt.Sprintf(formate, a...),
				},
			})
		}
		return WithStack(&withMessage{
			cause: err,
			err:   errors.Errorf(formate, a...),
		})
	default:
		return WithStack(&withMessage{
			cause: err,
			err:   errors.Errorf(formate, a...),
		})
	}

}

func (e *_error) UnmarshalJSON(b []byte) error {

	err := &ErrorView{}
	_err := json.Unmarshal(b, err)
	if _err != nil {
		return NewError(ErrorInternalError, _err.Error())
	}

	e.category = err.Category
	e.code = err.Code
	e.message = err.Message
	e.top = new(int)

	return nil
}

func (e *_error) MarshalJSON() ([]byte, error) {

	err := &ErrorView{
		Message:  e.message,
		Code:     e.code,
		Category: e.category,
	}

	b, _err := json.Marshal(err)
	if _err != nil {
		return nil, NewError(ErrorInternalError, _err.Error())
	}

	return b, nil
}

func Is(err error, target error) bool {
	causeErr, ok := errors.Cause(err).(*_error)
	causeTarget, _ok := errors.Cause(target).(*_error)
	if ok && _ok {
		return causeErr.code == causeTarget.code
	}

	return errors.Is(err, target)
}

func CompareWithBaseError(e baseError, target error) bool {
	causeTarget, ok := errors.Cause(target).(*_error)
	if !ok {
		return false
	}

	return e().code == causeTarget.code
}

func ConvertPostgresError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewError(ErrorResourceNotFound, err.Error())
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewError(ErrorConflict, pgErr.Error())
		}
	}

	return NewError(ErrorInternalError, err.Error())
}
