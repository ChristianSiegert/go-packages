package loggers

import (
	"context"
	"errors"
	"log"
)

type contextKey int

// key is used to identify the value carried by context.
const key contextKey = 0

// ErrNoLogger is the error returned by FromContext when context does not carry
// a logger.
var ErrNoLogger = errors.New("logger not carried by context")

// NewContext returns a new context that carries logger.
func NewContext(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

// FromContext returns the logger carried by ctx. If no logger is carried, error
// is ErrNoLogger.
func FromContext(ctx context.Context) (*log.Logger, error) {
	if ctx == nil {
		return nil, ErrNoLogger
	}

	logger, ok := ctx.Value(key).(*log.Logger)
	if !ok {
		return nil, ErrNoLogger
	}
	return logger, nil
}
