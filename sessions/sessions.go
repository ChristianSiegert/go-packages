// Package sessions provides session handling.
package sessions

import (
	"net/http"
	"time"
)

type Session struct {
	DateCreated time.Time
	Flashes     []Flash
	id          string
	isStored    bool
	store       Store
	UserId      string
	Values      map[interface{}]interface{}
}

// New returns a new session.
func New(store Store, id string) *Session {
	return &Session{
		DateCreated: time.Now(),
		id:          id,
		store:       store,
		Values:      make(map[interface{}]interface{}),
	}
}

// AddFlash adds a flash message to the session.
func (s *Session) AddFlash(message string, flashType ...string) {
	flash := Flash{
		Message: message,
	}

	if len(flashType) > 0 {
		flash.Type = flashType[0]
	}

	s.Flashes = append(s.Flashes, flash)
}

// Delete deletes the session from the session store.
func (s *Session) Delete(writer http.ResponseWriter) error {
	return s.store.Delete(writer, s)
}

// FlashAll returns all flashes and removes them from the session.
func (s *Session) FlashAll() []Flash {
	if len(s.Flashes) == 0 {
		return s.Flashes
	}

	flashes := s.Flashes
	s.Flashes = []Flash{}
	return flashes
}

// Id returns the session id.
func (s *Session) Id() string {
	return s.id
}

// RemoveFlash removes a flash message from the session.
func (s *Session) RemoveFlash(flash Flash) {
	for i := range s.Flashes {
		if s.Flashes[i] == flash {
			s.Flashes = append(s.Flashes[:i], s.Flashes[i+1:]...)
		}
	}
}

// Save saves the session to the session store.
func (s *Session) Save(writer http.ResponseWriter) error {
	return s.store.Save(writer, s)
}
