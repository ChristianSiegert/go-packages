// Package sqlitestores provides a session store backed by an SQLite database.
package sqlitestores

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/ChristianSiegert/go-packages/sessions"
	_ "github.com/mattn/go-sqlite3"
)

// Pattern for matching a session id.
var patternId = regexp.MustCompile("^[0-9a-zA-Z=/+]{1,}$")

type Store struct {
	cookieDomain string
	cookieName   string
	cookiePath   string
	db           *sql.DB

	// Duration after which sessions expire.
	Expiration time.Duration

	sessionStrength int
	tableName       string
}

// New returns a new SQLite session store. If a database table with the
// specified name does not exist, it is created.
func New(db *sql.DB, tableName, cookieName, cookieDomain, cookiePath string, strength int) (sessions.Store, error) {
	err := createSchema(db, tableName)
	if err != nil {
		return nil, err
	}

	return &Store{
		cookieDomain:    cookieDomain,
		cookieName:      cookieName,
		cookiePath:      cookiePath,
		db:              db,
		Expiration:      14 * 24 * time.Hour,
		sessionStrength: strength,
		tableName:       tableName,
	}, nil
}

// Delete deletes a session from the session store.
func (s *Store) Delete(writer http.ResponseWriter, session *sessions.Session) error {
	query := "DELETE FROM sessions WHERE id = ?"
	if _, err := s.db.Exec(query, session.Id()); err != nil {
		return err
	}

	s.deleteCookie(writer)
	session = nil
	return nil
}

// Save retrieves a session from the session store.
func (s *Store) Get(writer http.ResponseWriter, request *http.Request) (*sessions.Session, error) {
	cookie, err := request.Cookie(s.cookieName)

	if err == http.ErrNoCookie {
		return s.newSession()
	} else if err != nil {
		return nil, err
	}

	if !isId(cookie.Value) {
		s.deleteCookie(writer)
		return s.newSession()
	}

	session := sessions.New(s, cookie.Value)

	query := `
		SELECT
			data,
			dateCreated,
			flashes,
			userId
		FROM
			%s
		WHERE
			id = ?
		LIMIT 1
	`

	query = fmt.Sprintf(query, s.tableName)
	row := s.db.QueryRow(query, session.Id())
	var flashes []byte
	var values []byte

	err = row.Scan(
		&values,
		&session.DateCreated,
		&flashes,
		&session.UserId,
	)
	if err == sql.ErrNoRows {
		return s.newSession()
	}
	if err != nil {
		return nil, err
	}

	// Decode flashes
	encodedFlashes := bytes.NewBuffer(flashes)
	decoder := gob.NewDecoder(encodedFlashes)

	if err := decoder.Decode(&session.Flashes); err != nil {
		return nil, err
	}

	// Decode values
	encodedValues := bytes.NewBuffer(values)
	decoder = gob.NewDecoder(encodedValues)

	if err := decoder.Decode(&session.Values); err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		s.deleteCookie(writer)
		return s.newSession()
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

// Save saves a session in the session store.
func (s *Store) Save(writer http.ResponseWriter, session *sessions.Session) error {
	s.saveCookie(writer, session)

	query := `
		INSERT OR REPLACE INTO %s (
			data, dateCreated, flashes, id, userId
		) VALUES (
			?, ?, ?, ?, ?
		);
	`

	query = fmt.Sprintf(query, s.tableName)

	// Encode flashes
	flashesEncoded := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(flashesEncoded)
	if err := encoder.Encode(session.Flashes); err != nil {
		return err
	}

	// Encode values
	valuesEncoded := bytes.NewBuffer([]byte{})
	encoder = gob.NewEncoder(valuesEncoded)
	if err := encoder.Encode(session.Values); err != nil {
		return err
	}

	_, err := s.db.Exec(
		query,
		valuesEncoded.Bytes(),
		session.DateCreated,
		flashesEncoded.Bytes(),
		session.Id(),
		session.UserId,
	)
	return err
}

// newSession returns a new session with a randomly generated id.
func (s *Store) newSession() (*sessions.Session, error) {
	id, err := generateId(s.sessionStrength)
	if err != nil {
		return nil, err
	}
	return sessions.New(s, id), nil
}

func (s *Store) saveCookie(writer http.ResponseWriter, session *sessions.Session) {
	dateExpires := session.DateCreated.Add(s.Expiration)

	http.SetCookie(writer, &http.Cookie{
		Domain:   s.cookieDomain,
		Expires:  dateExpires,
		HttpOnly: true,
		MaxAge:   int(dateExpires.Sub(time.Now()).Seconds()),
		Name:     s.cookieName,
		Path:     s.cookiePath,
		Value:    session.Id(),
	})
}

func (s *Store) deleteCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Domain:   s.cookieDomain,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Name:     s.cookieName,
		Path:     s.cookiePath,
	})
}

// generateId generates a session id and encodes it in Base64.
func generateId(strength int) (string, error) {
	id := make([]byte, strength)

	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(id), nil
}

// isId checks whether id is a valid session id.
func isId(id string) bool {
	return patternId.MatchString(id)
}

func createSchema(db *sql.DB, tableName string) error {
	query := `
		CREATE TABLE IF NOT EXISTS %s (
			data BLOB,
			dateCreated TIMESTAMP NOT NULL,
			flashes BLOB,
			id TEXT PRIMARY KEY,
			userId TEXT
		);

		CREATE INDEX IF NOT EXISTS sessionsByDateCreated ON %s (
			dateCreated
		);

		CREATE INDEX IF NOT EXISTS sessionsByUserIdDateCreated ON %s (
			userId,
			dateCreated
		);
	`

	query = fmt.Sprintf(query, tableName, tableName, tableName)
	_, err := db.Exec(query)
	return err
}
