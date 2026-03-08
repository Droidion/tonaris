package httpapi

import (
	"net/http"
	"sync"

	"backend/internal/apperr"

	"github.com/danielgtaylor/huma/v2"
)

const internalErrorCode apperr.Code = "internal.error"

var humaErrorFactoryOnce sync.Once

type Problem struct {
	Type   string              `json:"type,omitempty" format:"uri" default:"about:blank" example:"https://example.com/errors/example" doc:"A URI reference to human-readable documentation for the error."`
	Title  string              `json:"title,omitempty" example:"Bad Request" doc:"A short, human-readable summary of the problem type. This value should not change between occurrences of the error."`
	Status int                 `json:"status,omitempty" example:"400" doc:"HTTP status code"`
	Detail string              `json:"detail,omitempty" example:"Property foo is required but is missing." doc:"A human-readable explanation specific to this occurrence of the problem."`
	Code   string              `json:"code,omitempty" example:"config.invalid_port" doc:"A stable machine-readable error code."`
	Errors []*huma.ErrorDetail `json:"errors,omitempty" doc:"Optional list of individual error details"`
}

func configureErrorResponses() {
	humaErrorFactoryOnce.Do(func() {
		huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
			return newProblem(status, msg, errs...)
		}
		huma.NewErrorWithContext = func(_ huma.Context, status int, msg string, errs ...error) huma.StatusError {
			return newProblem(status, msg, errs...)
		}
	})
}

func newProblem(status int, msg string, errs ...error) *Problem {
	problem := &Problem{
		Status: status,
		Title:  http.StatusText(status),
		Detail: msg,
		Code:   codeForStatus(status),
	}

	appErr, hasAppErr := firstAppError(errs...)
	if hasAppErr {
		problem.Status = apperr.Status(appErr)
		problem.Title = http.StatusText(problem.Status)
		problem.Detail = appErr.Message
		problem.Code = string(appErr.Code)
	}

	if !hasAppErr && problem.Status >= http.StatusInternalServerError {
		problem.Detail = "unexpected error occurred"
		problem.Code = string(internalErrorCode)
	} else if problem.Detail == "" {
		problem.Detail = msg
	}

	problem.Errors = collectProblemDetails(problem.Status, errs...)

	return problem
}

func collectProblemDetails(status int, errs ...error) []*huma.ErrorDetail {
	details := make([]*huma.ErrorDetail, 0, len(errs))

	for _, err := range errs {
		if err == nil {
			continue
		}

		if _, ok := apperr.As(err); ok {
			continue
		}

		if detailer, ok := err.(huma.ErrorDetailer); ok {
			details = append(details, detailer.ErrorDetail())
			continue
		}

		if status < http.StatusInternalServerError {
			details = append(details, &huma.ErrorDetail{Message: err.Error()})
		}
	}

	if len(details) == 0 {
		return nil
	}

	return details
}

func firstAppError(errs ...error) (*apperr.Error, bool) {
	for _, err := range errs {
		if appErr, ok := apperr.As(err); ok {
			return appErr, true
		}
	}

	return nil, false
}

func codeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "http.bad_request"
	case http.StatusUnauthorized:
		return "http.unauthorized"
	case http.StatusForbidden:
		return "http.forbidden"
	case http.StatusNotFound:
		return "http.not_found"
	case http.StatusMethodNotAllowed:
		return "http.method_not_allowed"
	case http.StatusConflict:
		return "http.conflict"
	case http.StatusUnprocessableEntity:
		return "http.unprocessable_entity"
	case http.StatusTooManyRequests:
		return "http.too_many_requests"
	default:
		if status >= http.StatusInternalServerError {
			return string(internalErrorCode)
		}
		if status == 0 {
			return ""
		}
		return "http.error"
	}
}

func (p *Problem) Error() string {
	if p.Detail != "" {
		return p.Detail
	}

	return p.Title
}

func (p *Problem) GetStatus() int {
	return p.Status
}

func (p *Problem) ContentType(ct string) string {
	if ct == "application/json" {
		return "application/problem+json"
	}
	if ct == "application/cbor" {
		return "application/problem+cbor"
	}
	return ct
}

var (
	_ error                  = (*Problem)(nil)
	_ huma.StatusError       = (*Problem)(nil)
	_ huma.ContentTypeFilter = (*Problem)(nil)
)
