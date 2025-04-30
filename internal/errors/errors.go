// Package errors provides common error definitions
package errors

import "errors"

// Common error definitions shared across packages
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("not found")

	// ErrDuplicate is returned when a duplicate entry is detected
	ErrDuplicate = errors.New("duplicate entry")

	// ErrInvalidInput is returned when the request contains invalid data
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")
)
