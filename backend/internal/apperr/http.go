package apperr

import "net/http"

func Status(err error) int {
	appErr, ok := As(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch appErr.Kind {
	case InvalidArgument:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case Unavailable:
		return http.StatusServiceUnavailable
	case Internal:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
