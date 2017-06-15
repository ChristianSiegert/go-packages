package sessions

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFromContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := NewSession(nil, "session123")
		newCtx := NewContext(ctx, session)

		if session2, isPresent := FromContext(newCtx); !isPresent {
			t.Error("Expected session to be present in context.")
		} else if !reflect.DeepEqual(session, session2) {
			t.Errorf("Expect session %#v, got %#v", session, session2)
		}
	}))
	defer server.Close()
}
