// Package sqlsessionstores provides a session store backed by an SQL
// database. Supported dialects are PostgreSQL and SQLite.
//
// You have to import the appropriate SQL driver yourself, e.g.:
//     _ "github.com/lib/pq"           // for PostgreSQL, or:
//     _ "github.com/mattn/go-sqlite3" // for SQLite
package sqlsessionstores

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
)

// Pattern is the pattern used to match a session ID.
var pattern = regexp.MustCompile("^[0-9a-zA-Z=/+]+$")

// KeyUserID is the key used to retrieve the user ID from session.Values and
// store it in an indexed table column. This makes it possible to delete all
// sessions of a particular user.
var KeyUserID = "user.id"

// authMethod is the method used to pass session IDs between server and client.
type authMethod string

const (
	// AuthMethodCookie means the session ID is passed via cookie.
	AuthMethodCookie = authMethod("cookie")

	// AuthMethodHeader means the session ID is passed via request header.
	AuthMethodHeader = authMethod("header")
)

// Dialect is the SQL dialect the store uses.
type dialect string

// Supported SQL dialects.
const (
	DialectPostgreSQL = dialect("postgres")
	DialectSQLite     = dialect("sqlite")
)

// Store contains information about the session store.
type Store struct {
	// Authentication options.
	AuthOptions AuthOptions

	// DB is the database in which the sessions table resides.
	DB *sql.DB

	// SQL dialect to use.
	Dialect dialect

	// Expiration is the duration after which sessions expire.
	Expiration time.Duration

	// Strength is the number of bytes to use for generating a session ID. The
	// higher the number, the more secure the session ID.
	Strength int

	// TableName is the name of the sessions table.
	TableName string
}

// AuthOptions is the authentification configuration for the store. If
// AuthMethod is AuthMethodCookie, Cookie… options are used. If AuthMethod is
//  AuthMethodHeader, Header… options are used.
type AuthOptions struct {
	AuthMethod authMethod

	// Cookie… fields are used when setting the cookie.
	CookieDomain string
	CookieName   string
	CookiePath   string

	// HeaderName is the name of the request header that is used to pass the
	// session ID.
	HeaderName string
}

// New returns a new Store. If a table with the specified name does not exist,
// it is created.
func New(db *sql.DB, tableName string, dialect dialect, authOptions AuthOptions) (*Store, error) {
	if err := createSchema(db, tableName, dialect); err != nil {
		return nil, err
	}

	store := &Store{
		AuthOptions: authOptions,
		DB:          db,
		Dialect:     dialect,
		Expiration:  14 * 24 * time.Hour,
		Strength:    40,
		TableName:   tableName,
	}

	return store, nil
}

func createSchema(db *sql.DB, tableName string, dialect dialect) error {
	query := fmt.Sprintf(
		queries[dialect][queryCreate],
		tableName,
		tableName,
		tableName,
		tableName,
		tableName,
	)

	_, err := db.Exec(query)
	return err
}

// Delete deletes a session from the store.
func (s *Store) Delete(writer http.ResponseWriter, sessionID string) error {
	query := fmt.Sprintf(queries[s.Dialect][queryDelete], s.TableName)
	if _, err := s.DB.Exec(query, sessionID); err != nil {
		return err
	}

	if s.AuthOptions.AuthMethod == AuthMethodCookie {
		s.deleteCookie(writer)
	}
	return nil
}

// DeleteMulti deletes sessions from the store that match the criteria specified
// in filter.
func (s *Store) DeleteMulti(filter *sessions.Filter) error {
	if filter != nil {
		return errors.New("filter not implemented")
	}

	query := "DELETE FROM %s"
	query = fmt.Sprintf(query, s.TableName)

	_, err := s.DB.Exec(query)
	return err
}

// Get gets a session from the store using the session ID stored in the session
// cookie.
func (s *Store) Get(writer http.ResponseWriter, request *http.Request) (sessions.Session, error) {
	var sessionID string

	switch s.AuthOptions.AuthMethod {
	case AuthMethodCookie:
		cookie, err := request.Cookie(s.AuthOptions.CookieName)

		if err == http.ErrNoCookie {
			return s.newSession()
		} else if err != nil {
			return nil, err
		} else if !isID(cookie.Value) {
			s.deleteCookie(writer)
			return s.newSession()
		}
		sessionID = cookie.Value
	case AuthMethodHeader:
		sessionID = request.Header.Get(s.AuthOptions.HeaderName)
	}

	if !isID(sessionID) {
		return s.newSession()
	}

	session := sessions.NewSession(s, sessionID)

	temp := struct {
		dateCreated    time.Time
		encodedFlashes []byte
		encodedValues  []byte
		flashes        []sessions.Flash
		userID         string
		values         map[string]string
	}{}

	query := fmt.Sprintf(queries[s.Dialect][queryGet], s.TableName)
	row := s.DB.QueryRow(query, session.ID())

	err := row.Scan(
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

	session.SetDateCreated(temp.dateCreated)
	session.SetIsStored(true)

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
// filter.
func (s *Store) GetMulti(filter *sessions.Filter) ([]sessions.Session, error) {
	return nil, errors.New("method not implemented")
}

// Save saves a session to the store and creates / updates the session cookie.
func (s *Store) Save(writer http.ResponseWriter, session sessions.Session) error {
	if s.AuthOptions.AuthMethod == AuthMethodCookie {
		s.saveCookie(writer, session)
	}

	query := fmt.Sprintf(queries[s.Dialect][querySave], s.TableName)

	encodedFlashes, err := json.Marshal(session.Flashes().GetAll())
	if err != nil {
		return err
	}

	encodedValues, err := json.Marshal(session.Values().GetAll())
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		query,
		encodedValues,
		session.DateCreated(),
		encodedFlashes,
		session.ID(),
		session.Values().Get(KeyUserID),
	)

	if err != nil {
		return err
	}

	session.SetIsStored(true)
	return nil
}

// SaveMulti saves the provided sessions.
func (s *Store) SaveMulti(sessions []sessions.Session) (e error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// If tx was not committed, rollback. If rollback fails, return rollback’s
	// error instead of the original error.
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			e = err
		}
	}()

	query := fmt.Sprintf(queries[s.Dialect][querySave], s.TableName)
	statement, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		encodedFlashes, err := json.Marshal(session.Flashes().GetAll())
		if err != nil {
			return err
		}

		encodedValues, err := json.Marshal(session.Values().GetAll())
		if err != nil {
			return err
		}

		_, err = statement.Exec(
			encodedValues,
			session.DateCreated(),
			encodedFlashes,
			session.ID(),
			session.Values().Get(KeyUserID),
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// newSession returns a new session with a randomly generated ID.
func (s *Store) newSession() (sessions.Session, error) {
	id, err := generateID(s.Strength)
	if err != nil {
		return nil, err
	}
	return sessions.NewSession(s, id), nil
}

func (s *Store) saveCookie(writer http.ResponseWriter, session sessions.Session) {
	dateExpires := session.DateCreated().Add(s.Expiration)

	http.SetCookie(writer, &http.Cookie{
		Domain:   s.AuthOptions.CookieDomain,
		Expires:  dateExpires,
		HttpOnly: true,
		MaxAge:   int(dateExpires.Sub(time.Now()).Seconds()),
		Name:     s.AuthOptions.CookieName,
		Path:     s.AuthOptions.CookiePath,
		Value:    session.ID(),
	})
}

func (s *Store) deleteCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Domain:   s.AuthOptions.CookieDomain,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Name:     s.AuthOptions.CookieName,
		Path:     s.AuthOptions.CookiePath,
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
	return pattern.MatchString(id)
}
