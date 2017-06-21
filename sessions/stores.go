package sessions

import (
	"net/http"
	"time"
)

// Store represents a session store.
type Store interface {
	// Delete deletes a session from the store, and deletes the session cookie.
	Delete(writer http.ResponseWriter, sessionID string) error

	// DeleteMulti deletes sessions from the store that match the criteria
	// specified in filter.
	DeleteMulti(filter *Filter) error

	// Get gets a session from the store using the session ID stored in the
	// session cookie.
	Get(http.ResponseWriter, *http.Request) (Session, error)

	// GetMulti gets sessions from the store that match the criteria specified
	// in filter.
	GetMulti(filter *Filter) ([]Session, error)

	// Save saves a session to the store and creates / updates the session
	// cookie.
	Save(http.ResponseWriter, Session) error

	// SaveMulti saves the provided sessions.
	SaveMulti([]Session) error
}

// Filter is used to limit DeleteMulti and GetMulti to sessions that match the
// criteria. Sessions match when 1) they have an ID or userID that is specified
// in IDs or UserIDs, and 2) their DateCreated is before DateCreatedBefore or
// after DateCreatedAfter. If both IDs and UserIDs are empty, sessions match
// regardless of their ID and session ID. If both DateCreatedBefore and
// DateCreatedAfter are zero, sessions match regardless of their DateCreated.
// Thus, with no filter set, all sessions match.
type Filter struct {
	DateCreatedAfter  time.Time
	DateCreatedBefore time.Time
	IDs               []string
	UserIDs           []string
}
