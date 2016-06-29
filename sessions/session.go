// Package sessions provides database sessions for PostgreSQL. sessions expects
// Db to be set before
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
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ChristianSiegert/go-packages/chttp"
	"golang.org/x/net/context"
)

const (
	cookieName = "s"
	cookiePath = "/"
)

// Database instance to use for saving and retrieving sessions.
var Db *sql.DB

// Duration in nanoseconds until a session expires.
var ExpirationLength = 14 * 24 * time.Hour

var regExpSessionId = regexp.MustCompile("[0-9a-zA-Z=/+]{88}")

type Session struct {
	dateCreated    time.Time
	flashes        []Flash
	id             string
	isAdmin        bool
	request        *http.Request
	responseWriter http.ResponseWriter
	userId         int64
	username       string

	// isPersistent indicates whether the session was retrieved from the
	// database. If false, the session was not yet saved to the database.
	isPersistent bool
}

func newSession(ctx context.Context) (*Session, error) {
	responseWriter, request, ok := chttp.FromContext(ctx)
	if !ok {
		return nil, errors.New("sessions.newSession: http.ResponseWriter and http.Request are not provided by ctx.")
	}

	sessionId, err := generateSessionId()
	if err != nil {
		return nil, err
	}

	return &Session{
		dateCreated:    time.Now(),
		id:             sessionId,
		request:        request,
		responseWriter: responseWriter,
	}, nil
}

// Id returns the session id.
func (s *Session) Id() string {
	return s.id
}

// IsAdmin returns whether the user has admin rights.
func (s *Session) IsAdmin() bool {
	return s.isAdmin
}

// UserId returns the id of the user.
func (s *Session) UserId() int64 {
	return s.userId
}

// SetUserId sets the id of the user.
func (s *Session) SetUserId(userId int64) error {
	if s.userId != 0 {
		return fmt.Errorf("User id is already set.")
	}
	s.userId = userId
	return nil
}

func Get(ctx context.Context) (*Session, error) {
	responseWriter, request, ok := chttp.FromContext(ctx)
	if !ok {
		return nil, errors.New("sessions.Get: http.ResponseWriter and http.Request are not provided by ctx.")
	}

	cookie, err := request.Cookie(cookieName)

	if err == http.ErrNoCookie {
		return newSession(ctx)
	}

	if !isSessionId(cookie.Value) {
		expireCookie(responseWriter, cookie)
		return newSession(ctx)
	}

	session := &Session{
		id:             cookie.Value,
		request:        request,
		responseWriter: responseWriter,
	}

	query := "SELECT date_created, user_id FROM sessions WHERE id = $1 LIMIT 1"
	row := Db.QueryRow(query, session.id)

	err = row.Scan(
		&session.dateCreated,
		&session.userId,
	)

	if err == sql.ErrNoRows {
		expireCookie(responseWriter, cookie)
		return newSession(ctx)
	} else if err != nil {
		return nil, err
	}

	session.isPersistent = true
	return session, nil
}

func (s *Session) Save() error {
	if _, err := s.request.Cookie(cookieName); err == http.ErrNoCookie {
		dateExpires := s.dateCreated.Add(ExpirationLength)

		http.SetCookie(s.responseWriter, &http.Cookie{
			Expires:  dateExpires,
			HttpOnly: true,
			MaxAge:   int(dateExpires.Sub(time.Now()).Seconds()),
			Name:     cookieName,
			Path:     cookiePath,
			Value:    s.id,
		})
	}

	query := `
		INSERT INTO sessions (
			date_created, id, user_id
		) VALUES (
			$1, $2, $3
		) ON CONFLICT DO NOTHING
	`

	_, err := Db.Exec(query, s.dateCreated, s.id, s.userId)
	return err
}

// Expire expires the session by deleting it from the database, and by deleting
// the client’s session cookie.
func (s *Session) Expire() error {
	query := "DELETE FROM sessions WHERE id = $1"
	if _, err := Db.Exec(query, s.id); err != nil {
		return err
	}

	s.userId = 0
	s.isPersistent = false

	// If a cookie exists, expire it.
	if cookie, err := s.request.Cookie(cookieName); err == nil {
		expireCookie(s.responseWriter, cookie)
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

// IsSignedIn returns whether the user is signed in.
func (s *Session) IsSignedIn() bool {
	return s.userId != 0
}

func (s *Session) SignIn(userId int64, username string, isAdmin bool) {
	s.isAdmin = isAdmin
	s.userId = userId
	s.username = username
}

// AddFlashError adds a flash of type error to the session.
func (s *Session) AddFlashError(message string) Flash {
	flash := Flash{
		Message: message,
		Type:    FlashTypeError,
	}

	s.flashes = append(s.flashes, flash)
	return flash
}

// AddFlashInfo adds a flash of type info to the session.
func (s *Session) AddFlashInfo(message string) Flash {
	flash := Flash{
		Message: message,
		Type:    FlashTypeInfo,
	}

	s.flashes = append(s.flashes, flash)
	return flash
}

// FlashAll returns all flashes belonging to session. By doing so, they are
// removed from session.
func (s *Session) FlashAll() []Flash {
	f := s.flashes

	if len(f) > 0 {
		s.flashes = []Flash{}

		if s.isPersistent {
			if err := s.Save(); err != nil {
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

// generateSessionId generates a 528-bit long session token and encodes it in
// Base64.
func generateSessionId() (string, error) {
	length := 66
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

// contextKey is used for attaching Session to context.Context.
type contextKey int

// key used for storing and retrieving the session from context.Context.
const key contextKey = 0

// NewContext returns a new Context carrying session.
func NewContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, key, session)
}

// FromContext extracts the session from ctx, if present.
func FromContext(ctx context.Context) (*Session, bool) {
	if ctx == nil {
		return nil, false
	}
	session, ok := ctx.Value(key).(*Session)
	return session, ok
}
