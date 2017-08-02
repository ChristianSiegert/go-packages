package sessions

import (
	"context"
	"errors"
)

type contextKey int

// key is used to identify the value carried by context.
const key contextKey = 0

// ErrNoSession is the error returned by FromContext when context does not carry
// a session.
var ErrNoSession = errors.New("session not carried by context")

// NewContext returns a new context that carries session.
func NewContext(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, key, session)
}

// FromContext returns the session carried by ctx. If no session is carried,
// error is ErrNoSession.
func FromContext(ctx context.Context) (Session, error) {
	if ctx == nil {
		return nil, ErrNoSession
	}

	session, ok := ctx.Value(key).(Session)
	if !ok {
		return nil, ErrNoSession
	}
	return session, nil
}
