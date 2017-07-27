package buffer

import "errors"

// Error returned by package buffer
var (
	// ErrClosedBuffer stands for accessing a closed buffer
	ErrClosedBuffer = errors.New("closed buffer")

	// ErrInvalidBufferSize occurs when init a buffer with invalid size parameter
	ErrInvalidBufferSize = errors.New("invalid buffer size")

	// ErrClosedBuffer stands for accessing a closed pool
	ErrClosedPool = errors.New("closed pool")
	
	// ErrInvalidPoolSize occurs when init a pool with invalid size parameter
	ErrInvalidPoolSize = errors.New("invalid pool size")
)
