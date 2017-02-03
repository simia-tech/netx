package netx

import "errors"

var (
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNotSupported       = errors.New("not supported")
)
