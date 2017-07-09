// Package sessions provides HTTP(S) sessions.
package sessions

import (
	"net/http"
	"time"
)

// Session represents an HTTP(S) session.
type Session interface {
	// DateCreated returns the session’s creation date.
	DateCreated() time.Time

	// Delete deletes the session from the session store.
	Delete(http.ResponseWriter) error

	// Flashes returns the session’s flash container.
	Flashes() Flashes

	// ID returns the session’s ID.
	ID() string

	// Save saves the session to the session store.
	Save(http.ResponseWriter) error

	// SetDateCreated sets the session’s creation date.
	SetDateCreated(time.Time)

	// Store returns the session store.
	Store() Store

	// Values returns the session’s value container.
	Values() Values
}

// session is an unexported type that implements the Session interface.
type session struct {
	dateCreated time.Time
	flashes     Flashes
	id          string
	store       Store
	values      Values
}

// NewSession returns a new session. The session has not been saved to the
// session store yet. To do that, call Save.
func NewSession(store Store, id string) Session {
	return &session{
		dateCreated: time.Now(),
		flashes:     NewFlashes(),
		id:          id,
		store:       store,
		values:      NewValues(),
	}
}

// DateCreated returns the session’s creation date.
func (s *session) DateCreated() time.Time {
	return s.dateCreated
}

// Delete deletes the session from the session store.
func (s *session) Delete(writer http.ResponseWriter) error {
	return s.store.Delete(writer, s.ID())
}

// Flashes returns the session’s flash container.
func (s session) Flashes() Flashes {
	return s.flashes
}

// ID returns the session’s id.
func (s *session) ID() string {
	return s.id
}

// Save saves the session to the session store.
func (s *session) Save(writer http.ResponseWriter) error {
	return s.store.Save(writer, s)
}

// SetDateCreated sets the session’s creation date.
func (s *session) SetDateCreated(date time.Time) {
	s.dateCreated = date
}

// Store returns the session store.
func (s session) Store() Store {
	return s.store
}

// Values returns the session’s value container.
func (s session) Values() Values {
	return s.values
}
