package sessions

import "context"

// contextKey is used for attaching a session to a context.
type contextKey int

// key used for storing and retrieving the session from a context.
const key contextKey = 0

// NewContext returns a new context that carries session.
func NewContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, key, session)
}

// FromContext extracts the session from ctx. If no session is carried by ctx,
// the second return argument is false.
func FromContext(ctx context.Context) (*Session, bool) {
	if ctx == nil {
		return nil, false
	}
	session, ok := ctx.Value(key).(*Session)
	return session, ok
}
