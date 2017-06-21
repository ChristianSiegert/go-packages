package sqlitesessionstores

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/ChristianSiegert/go-packages/sessions"
)

var dateCreated = time.Date(2099, 12, 31, 13, 14, 15, 16, &time.Location{})

func setUp() (*sql.DB, sessions.Store, error) {
	filename := path.Join(os.TempDir(), "test.sqlite")

	// Clean up
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("Removing database file failed: %s", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Opening database failed: %s", err)
	}

	// Create store instance
	store, err := New(db, "test_sessions", "s", "", "", 10)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("Creating store failed:%s", err)
	}
	return db, store, nil
}

func tearDown(db *sql.DB) {
	db.Close()
}

func Test(t *testing.T) {
	db, store, err := setUp()
	if err != nil {
		t.Error(err)
	}
	defer tearDown(db)

	// Create routes
	mux := http.NewServeMux()
	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		testSave(w, r, t, store)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		testGet(w, r, t, store)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		testDelete(w, r, t, store)
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("Creating cookie jar failed: %s", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Serve pages
	server := httptest.NewServer(mux)
	defer server.Close()

	if _, err := client.Get(server.URL + "/save"); err != nil {
		t.Errorf("GET request failed: %s", err)
	} else if _, err := client.Get(server.URL + "/get"); err != nil {
		t.Errorf("GET request failed: %s", err)
	} else if _, err := client.Get(server.URL + "/delete"); err != nil {
		t.Errorf("GET request failed: %s", err)
	}
}

func testSave(writer http.ResponseWriter, request *http.Request, t *testing.T, store sessions.Store) {
	session := sessions.NewSession(store, "session123")
	session.SetDateCreated(dateCreated)
	session.Flashes().AddNew("lorem ipsum", "info")
	session.Values().Set("user.id", "user1")

	if err := store.Save(writer, session); err != nil {
		t.Errorf("Saving session failed: %s", err)
	} else if writer.Header().Get("Set-Cookie") == "" {
		t.Errorf("Expected header Set-Cookie to be set.")
	}
}

func testGet(writer http.ResponseWriter, request *http.Request, t *testing.T, store sessions.Store) {
	expectedSession := sessions.NewSession(store, "session123")
	expectedSession.SetDateCreated(dateCreated)
	expectedSession.Flashes().AddNew("lorem ipsum", "info")
	expectedSession.Values().Set("user.id", "user1")

	session, err := store.Get(writer, request)
	if err != nil {
		t.Errorf("Getting session failed: %s", err)
	} else if !session.DateCreated().Equal(expectedSession.DateCreated()) {
		t.Errorf("Expected DateCreated %q, got %q.", session.DateCreated(), expectedSession.DateCreated())
	} else if !reflect.DeepEqual(session.Flashes(), expectedSession.Flashes()) {
		t.Errorf("Expected Flashes %#v, got %#v", expectedSession.Flashes(), session.Flashes())
	} else if session.ID() != expectedSession.ID() {
		t.Errorf("Expected ID %q, got %q.", expectedSession.ID(), session.ID())
	} else if !reflect.DeepEqual(session.Values(), expectedSession.Values()) {
		t.Errorf("Expected Values %#v, got %#v", expectedSession.Values(), session.Values())
	}
}

func testDelete(writer http.ResponseWriter, request *http.Request, t *testing.T, store sessions.Store) {
	if err := store.Delete(writer, "session123"); err != nil {
		t.Errorf("Deleting session failed: %s", err)
	}

	if session, err := store.Get(writer, request); err != nil {
		t.Errorf("Getting session failed: %s", err)
	} else if session.ID() == "session123" {
		t.Errorf("Expected random session ID, got old session ID %q.", session.ID())
	}
}

func TestMulti(t *testing.T) {
	db, store, err := setUp()
	if err != nil {
		t.Error(err)
	}
	defer tearDown(db)

	sessionA := sessions.NewSession(store, "a")
	sessionA.Flashes().AddNew("lorem", "ipsum")
	sessionA.SetDateCreated(time.Date(2090, 11, 10, 9, 8, 7, 6, &time.Location{}))
	sessionA.Values().Set(KeyUserID, "user-a")

	ss := []sessions.Session{
		sessionA,
		sessions.NewSession(store, "b"),
		sessions.NewSession(store, "c"),
	}

	if err := store.SaveMulti(ss); err != nil {
		t.Errorf("SaveMulti failed: %s", err)
	}

	ss2, err := store.GetMulti(nil)
	if err != nil {
		t.Errorf("GetMulti failed: %s", err)
	} else if !reflect.DeepEqual(ss2, ss) {
		t.Errorf("Expected sessions %#v, got %#v", ss, ss2)
	}

	if err := store.DeleteMulti(nil); err != nil {
		t.Errorf("DeleteMulti failed: %s", err)
	}
	if ss3, err := store.GetMulti(nil); err != nil {
		t.Errorf("Getting sessions failed: %s", err)
	} else if len(ss3) != 0 {
		t.Errorf("Expected 0 sessions, got %d.", len(ss3))
	}
}
