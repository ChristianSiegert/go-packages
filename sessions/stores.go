package sessions

import "net/http"

type Store interface {
	// Delete session.
	Delete(http.ResponseWriter, *Session) error

	// Get session.
	Get(http.ResponseWriter, *http.Request) (*Session, error)

	// Save session.
	Save(http.ResponseWriter, *Session) error
}
