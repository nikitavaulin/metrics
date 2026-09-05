package commonerrors

import "errors"

var (
	ErrNotFound error = errors.New("not found")
)
