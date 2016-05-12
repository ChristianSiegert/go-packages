package sessions

import (
	"reflect"
	"testing"
)

func TestGenerateSessionId(t *testing.T) {
	sessionId, err := generateSessionId()

	if err != nil {
		t.Errorf("Couldn’t generate session id: %s", err)
		return
	}

	if !regExpSessionId.MatchString(sessionId) {
		t.Errorf("Generated session id is not valid: %s", sessionId)
	}
}

func TestIsSessionId(t *testing.T) {
	validSessionId := "0AoDTgpn3KjoWrrA9fR61PpWvm9wx9h46DCpkhF3d0+tNRFpj5goHTt4S3Hls7+BlA1ujy98c6erYbBPYUxqyQ=="
	invalidSessionId := "tNRFpj5goHTt4S3Hls7+BlA1ujy98c6erYbBPYUxqyQ=="

	if !isSessionId(validSessionId) {
		t.Errorf("Valid session id was not recognized as valid: %s", validSessionId)
	}

	if isSessionId(invalidSessionId) {
		t.Errorf("Invalid session id was not recognized as invalid: %s", invalidSessionId)
	}
}

func TestSession_RemoveFlash(t *testing.T) {
	session := &Session{}

	var (
		flashError1 = session.AddFlashError("foo error 1")
		flashError2 = session.AddFlashError("foo error 2")
		flashError3 = session.AddFlashError("foo error 3")
		flashInfo   = session.AddFlashInfo("foo info")
		// flashWarning = session.AddFlashWarnMessage("foo warning")
	)

	// We expect all 5 flashes to have been added
	expected1 := []Flash{
		flashError1,
		flashError2,
		flashError3,
		flashInfo,
		// flashWarning,
	}

	if !reflect.DeepEqual(session.flashes, expected1) {
		t.Fatalf("session.Flashes is %q, expected %q.", session.flashes, expected1)
	}

	// We expect flashError2 to have been removed
	session.RemoveFlash(flashError2)

	expected2 := []Flash{
		flashError1,
		flashError3,
		flashInfo,
		// flashWarning,
	}

	if !reflect.DeepEqual(session.flashes, expected2) {
		t.Fatalf("session.Flashes is %q, expected %q.", session.flashes, expected2)
	}

	// Removing unknown flash shouldn’t remove anything
	session.RemoveFlash(Flash{
		Message: "custom flash",
		Type:    123,
	})

	if !reflect.DeepEqual(session.flashes, expected2) {
		t.Fatalf("session.Flashes is %q, expected %q.", session.flashes, expected2)
	}
}
