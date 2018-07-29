package exception

import "errors"

var (
	// ErrDataCorrupted is a data corrupted exception where manual data has affected enums.
	ErrDataCorrupted = errors.New("Data Corrupted Exception")
	// ErrDatabase is a database exception.
	ErrDatabase = errors.New("Database Exception")
	// ErrInvalidInput is an invalid input exception.
	ErrInvalidInput = errors.New("Invalid Input Exception")
	// ErrVersionInvalid is an optimistic locking exception
	// that is thrown when the record has been updated.
	ErrVersionInvalid = errors.New("Invalid Version Exception")
	// ErrUnauthorized is an unauthorized exception.
	ErrUnauthorized = errors.New("Unauthorized Exception")
	// ErrNotFound is a not found exception.
	ErrNotFound = errors.New("Not Found Exception")
	// ErrInvalidResourceState is an invalid resource state.
	ErrInvalidResourceState = errors.New("Invalid Resource State")
	// ErrServiceUnavailable is a service unavailable exception.
	ErrServiceUnavailable = errors.New("Service Unavailable Exception")
	// ErrUnhandledException is an unhandled exception.
	ErrUnhandledException = errors.New("Unhandled Exception")
)
