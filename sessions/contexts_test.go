package sessions

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFromContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		session := NewSession(nil, "session123")
		ctx := NewContext(request.Context(), session)

		if session2, err := FromContext(ctx); err != nil {
			t.Error("Expected session to be carried by context.")
		} else if !reflect.DeepEqual(session, session2) {
			t.Errorf("Expect session %#v, got %#v", session, session2)
		}
	}))
	defer server.Close()
}
