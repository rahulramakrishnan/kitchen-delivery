package exception

import "errors"

var (
	// ErrDataCorrupted is a data corrupted exception where
	// manual data has affected enums.
	ErrDataCorrupted = errors.New("Data Corrupted Exception")
	// ErrDatabase is a database exception - retriable.
	ErrDatabase = errors.New("Database Exception")
	// ErrInvalidInput is an invalid input exception - not retriable.
	ErrInvalidInput = errors.New("Invalid Input Exception")
	// ErrVersionInvalid is an optimistic locking exception
	// that is thrown when the record has been updated - retriable.
	ErrVersionInvalid = errors.New("Invalid Version Exception")
	// ErrUnauthorized is an unauthorized exception - not retriable.
	ErrUnauthorized = errors.New("Unauthorized Exception")
	// ErrNotFound is a not found exception - not retriable.
	ErrNotFound = errors.New("Not Found Exception")
	// ErrInvalidResourceState is an invalid resource state - not retriable.
	ErrInvalidResourceState = errors.New("Invalid Resource State")
	// ErrFullShelf is a full shelf - retriable.
	ErrFullShelf = errors.New("Shelf Space Is Full")
	// ErrServiceUnavailable is a service unavailable exception - retriable.
	ErrServiceUnavailable = errors.New("Service Unavailable Exception")
	// ErrUnhandledException is an unhandled exception - not retriable.
	ErrUnhandledException = errors.New("Unhandled Exception")
)
