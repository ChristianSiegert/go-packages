package sessions

import (
	"testing"
	"time"
)

func TestSession_DateCreated(t *testing.T) {
	session := NewSession(nil, "session123")
	now := time.Now()

	if session.DateCreated().Unix() < now.Unix()-1 || session.DateCreated().Unix() > now.Unix()+1 {
		t.Errorf("Unexpected DateCreated %s", session.DateCreated())
	}
}

func TestSession_ID(t *testing.T) {
	id := "session123"
	session := NewSession(nil, id)

	if session.ID() != id {
		t.Errorf("Expected ID %q, got %q.", id, session.ID())
	}
}

func TestSession_SetDateCreated(t *testing.T) {
	session := NewSession(nil, "session123")
	date := time.Unix(30, 0)
	session.SetDateCreated(date)

	if !session.DateCreated().Equal(date) {
		t.Errorf("Expected DateCreated %s, got %s", date, session.DateCreated())
	}
}
