// Package sessions provides database sessions for SQLite.
// An SQL table with the following structure must exist:
//		CREATE TABLE sessions (
//			dateCreated INTEGER NOT NULL,
//			id TEXT UNIQUE,
//			userId INTEGER NOT NULL
//		);
package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	cookieName = "s"
	cookiePath = "/"
)

const (
	FlashTypeError = iota
	FlashTypeInfo
)

// Database instance to use for saving and retrieving sessions.
var Db *sql.DB

// Duration in nanoseconds until a session expires.
var ExpirationLength = 14 * 24 * time.Hour

var regExpSessionId = regexp.MustCompile("[0-9a-zA-Z=/+]{88}")

type Session struct {
	dateCreated time.Time
	flashes     []Flash
	id          string
	userId      uint64

	// isPersistent indicates whether the session was retrieved from the
	// database. If false, the session was not yet saved to the database.
	isPersistent bool
}

func newSession() (*Session, error) {
	sessionId, err := generateSessionId()

	if err != nil {
		return nil, err
	}

	return &Session{
		dateCreated: time.Now(),
		id:          sessionId,
	}, nil
}

// Id returns the session id.
func (s *Session) Id() string {
	return s.id
}

// UserId returns the id of the user.
func (s *Session) UserId() uint64 {
	return s.userId
}

// SetUserId sets the id of the user.
func (s *Session) SetUserId(userId uint64) error {
	if s.userId != 0 {
		return fmt.Errorf("User id is already set.")
	}
	s.userId = userId
	return nil
}

func Get(responseWriter http.ResponseWriter, request *http.Request) (*Session, error) {
	cookie, err := request.Cookie(cookieName)

	if err == http.ErrNoCookie {
		return newSession()
	}

	if !isSessionId(cookie.Value) {
		expireCookie(responseWriter, cookie)
		return newSession()
	}

	session := &Session{
		id: cookie.Value,
	}

	query := "SELECT dateCreated, userId FROM sessions WHERE id = ? LIMIT 1"
	row := Db.QueryRow(query, session.id)

	var tempDateCreated int64

	err = row.Scan(
		&tempDateCreated,
		&session.userId,
	)

	if err == sql.ErrNoRows {
		expireCookie(responseWriter, cookie)
		return newSession()
	} else if err != nil {
		return nil, err
	}

	session.dateCreated = time.Unix(tempDateCreated, 0)
	session.isPersistent = true
	return session, nil
}

func (s *Session) Save(responseWriter http.ResponseWriter, request *http.Request) error {
	if _, err := request.Cookie(cookieName); err == http.ErrNoCookie {
		dateExpires := s.dateCreated.Add(ExpirationLength)

		http.SetCookie(responseWriter, &http.Cookie{
			Expires:  dateExpires,
			HttpOnly: true,
			MaxAge:   int(dateExpires.Sub(time.Now()).Seconds()),
			Name:     cookieName,
			Path:     cookiePath,
			Value:    s.id,
		})
	}

	query := `
		INSERT OR REPLACE INTO sessions (
			dateCreated, id, userId
		) VALUES (
			?, ?, ?
		)
	`

	_, err := Db.Exec(query, s.dateCreated.Unix(), s.id, s.userId)
	if err != nil {
		return err
	}
	return nil
}

// Expire expires the session by deleting it from the database, and by deleting
// the client’s session cookie.
func (s *Session) Expire(responseWriter http.ResponseWriter, request *http.Request) error {
	query := "DELETE FROM sessions WHERE id = ?"
	if _, err := Db.Exec(query, s.id); err != nil {
		return err
	}

	s.userId = 0
	s.isPersistent = false

	// If a cookie exists, expire it.
	if cookie, err := request.Cookie(cookieName); err == nil {
		expireCookie(responseWriter, cookie)
	}

	return nil
}

func expireCookie(responseWriter http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(-24 * time.Hour)
	cookie.MaxAge = -1
	cookie.Path = cookiePath
	cookie.Value = ""
	http.SetCookie(responseWriter, cookie)
}

func (s *Session) IsSignedIn() bool {
	return s.userId != 0
}

func (s *Session) AddFlashErrorMessage(message string) Flash {
	flash := Flash{
		Message: message,
		Type:    FlashTypeError,
	}

	s.flashes = append(s.flashes, flash)
	return flash
}

func (s *Session) AddFlashInfoMessage(message string) Flash {
	flash := Flash{
		Message: message,
		Type:    FlashTypeInfo,
	}

	s.flashes = append(s.flashes, flash)
	return flash
}

// FlashAll returns all flashes belonging to session. By doing so, they are
// removed from session.
func (s *Session) FlashAll(responseWriter http.ResponseWriter, request *http.Request) []Flash {
	f := s.flashes

	if len(f) > 0 {
		s.flashes = []Flash{}

		if s.isPersistent {
			if err := s.Save(responseWriter, request); err != nil {
				log.Printf("sessions.Session.FlashAll: %s", err)
			}
		}
	}

	return f
}

func (s *Session) RemoveFlash(flash Flash) {
	for i, v := range s.flashes {
		if flash == v {
			s.flashes = append(s.flashes[:i], s.flashes[i+1:]...)
		}
	}
}

// generateSessionId generates a unique identifier and encodes it in Base64 so
// it can be stored in cookies.
func generateSessionId() (string, error) {
	length := 64
	key := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// isSessionId checks whether id is a valid session id.
func isSessionId(id string) bool {
	return regExpSessionId.MatchString(id)
}

// DeleteExpired deletes sessions from the database that were created before
// maxDateCreated.
func DeleteExpired(maxDateCreated time.Time) error {
	query := "DELETE FROM sessions WHERE dateCreated < ?"
	_, err := Db.Exec(query, maxDateCreated.Unix())
	return err
}

// CronDeleteExpired calls DeleteExpired in regular intervals.
func CronDeleteExpired(interval time.Duration, expirationLength time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				maxDateCreated := time.Now().Add(-expirationLength)
				if err := DeleteExpired(maxDateCreated); err != nil {
					fmt.Errorf("sessions: Deleting expired sessions failed: %s", err)
				}
			}
		}
	}()
}
