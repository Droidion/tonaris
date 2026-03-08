package apperr

import "errors"

type Error struct {
	Kind    Kind
	Code    Code
	Message string
	Err     error
}

func New(kind Kind, code Code, message string) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, kind Kind, code Code, message string) *Error {
	if err == nil {
		return New(kind, code, message)
	}

	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	var appErr *Error
	if !errors.As(err, &appErr) {
		return nil, false
	}

	return appErr, true
}

func HasKind(err error, kind Kind) bool {
	appErr, ok := As(err)
	return ok && appErr.Kind == kind
}

func HasCode(err error, code Code) bool {
	appErr, ok := As(err)
	return ok && appErr.Code == code
}
