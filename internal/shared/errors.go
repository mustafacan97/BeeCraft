package shared

import "errors"

var (
	ErrInvalidContext = errors.New("invalid context value")
	ErrMissingContext = errors.New("missing context value")

	ErrNotFound = errors.New("resource not found")

	ErrUnauthorized = errors.New("unauthorized access")

	ErrAlreadyExists = errors.New("resource already exists")

	ErrValidation = errors.New("validation failed")

	ErrOperationFailed = errors.New("operation failed")

	ErrInternal = errors.New("internal server error")
)
