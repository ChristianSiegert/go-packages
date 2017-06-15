// Package sqlitesessionstores provides a session store backed by an SQLite
// database.
package sqlitesessionstores

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/ChristianSiegert/go-packages/sessions"

	// Register SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// Pattern for matching a session ID.
var patternID = regexp.MustCompile("^[0-9a-zA-Z=/+]+$")

// KeyUserID is used to retrieve the user ID from the session.Values container
// and store it in the table in an indexed column. This makes it possible to
// delete all sessions of a particular user.
var KeyUserID = "user.id"

// Store contains information about the session store.
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

// Delete deletes a session from the store, and deletes the session cookie.
func (s *Store) Delete(writer http.ResponseWriter, sessionID string) error {
	query := "DELETE FROM %s WHERE id = ?"
	query = fmt.Sprintf(query, s.tableName)
	if _, err := s.db.Exec(query, sessionID); err != nil {
		return err
	}

	s.deleteCookie(writer)
	return nil
}

// DeleteMulti deletes sessions from the store that match the criteria specified
// in options.
func (s *Store) DeleteMulti(options *sessions.StoreOptions) error {
	return errors.New("method no implemented")
}

// Get gets a session from the store using the session ID stored in the session
// cookie.
func (s *Store) Get(writer http.ResponseWriter, request *http.Request) (sessions.Session, error) {
	cookie, err := request.Cookie(s.cookieName)

	if err == http.ErrNoCookie {
		return s.newSession()
	} else if err != nil {
		return nil, err
	}

	if !isID(cookie.Value) {
		s.deleteCookie(writer)
		return s.newSession()
	}

	session := sessions.NewSession(s, cookie.Value)

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

	temp := struct {
		dateCreated    time.Time
		encodedFlashes []byte
		encodedValues  []byte
		flashes        []sessions.Flash
		userID         string
		values         map[string]string
	}{}

	query = fmt.Sprintf(query, s.tableName)
	row := s.db.QueryRow(query, session.ID())

	err = row.Scan(
		&temp.encodedValues,
		&temp.dateCreated,
		&temp.encodedFlashes,
		&temp.userID,
	)
	if err == sql.ErrNoRows {
		s.deleteCookie(writer)
		return s.newSession()
	} else if err != nil {
		return nil, err
	}

	// Date
	session.SetDateCreated(temp.dateCreated)

	// Decode flashes
	flashes, err := sessions.FlashesFromJSON(temp.encodedFlashes)
	if err != nil {
		return nil, err
	}
	session.Flashes().Add(flashes...)

	// Decode values
	values, err := sessions.ValuesFromJSON(temp.encodedValues)
	if err != nil {
		return nil, err
	}
	session.Values().SetAll(values)

	return session, nil
}

// GetMulti gets sessions from the store that match the criteria specified in
// options.
func (s *Store) GetMulti(options *sessions.StoreOptions) ([]sessions.Session, error) {
	return nil, errors.New("method no implemented")
}

// Save saves a session to the store and creates / updates the session cookie.
func (s *Store) Save(writer http.ResponseWriter, session sessions.Session) error {
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
	encodedFlashes, err := json.Marshal(session.Flashes().GetAll())
	if err != nil {
		return err
	}

	// Encode values
	encodedValues, err := json.Marshal(session.Values().GetAll())
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		query,
		encodedValues,
		session.DateCreated(),
		encodedFlashes,
		session.ID(),
		session.Values().Get(KeyUserID),
	)
	return err
}

// newSession returns a new session with a randomly generated ID.
func (s *Store) newSession() (sessions.Session, error) {
	id, err := generateID(s.sessionStrength)
	if err != nil {
		return nil, err
	}
	return sessions.NewSession(s, id), nil
}

func (s *Store) saveCookie(writer http.ResponseWriter, session sessions.Session) {
	dateExpires := session.DateCreated().Add(s.Expiration)

	http.SetCookie(writer, &http.Cookie{
		Domain:   s.cookieDomain,
		Expires:  dateExpires,
		HttpOnly: true,
		MaxAge:   int(dateExpires.Sub(time.Now()).Seconds()),
		Name:     s.cookieName,
		Path:     s.cookiePath,
		Value:    session.ID(),
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

// generateID generates a session ID and encodes it in Base64.
func generateID(strength int) (string, error) {
	id := make([]byte, strength)

	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(id), nil
}

// isID checks whether id is a valid session ID.
func isID(id string) bool {
	return patternID.MatchString(id)
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

// SaveMulti saves the provided sessions.
func (s *Store) SaveMulti(sessions []sessions.Session) error {
	return errors.New("method no implemented")
}
