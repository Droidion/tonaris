package apperr

type Kind string
type Code string

const (
	InvalidArgument Kind = "invalid_argument"
	NotFound        Kind = "not_found"
	Conflict        Kind = "conflict"
	Unauthorized    Kind = "unauthorized"
	Forbidden       Kind = "forbidden"
	Internal        Kind = "internal"
	Unavailable     Kind = "unavailable"
)
