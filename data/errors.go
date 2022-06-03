package data

import "errors"

// ErrNotFound should be returned when the resource does not exist.
var ErrNotFound = errors.New("resource not found")
