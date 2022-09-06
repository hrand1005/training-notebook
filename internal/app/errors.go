package app

import (
	"errors"
)

// ErrNotFound should be used for failures to obtain a resource by request parameters
var ErrNotFound = errors.New("failed to find the resource")

// ErrServiceFailure should be used for failures that are specific to a service's implementation
var ErrServiceFailure = errors.New("the service failed due to internal reasons")
