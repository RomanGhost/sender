package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sender/internal/server/web/handlers"
	"testing"
)

func TestKeysGenerateHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/keys/generate", nil)
	rr := httptest.NewRecorder()

	handlers.KeysGenerateHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected Content-Type %s, got %s", expectedContentType, rr.Header().Get("Content-Type"))
	}

	var keys handlers.Keys
	if err := json.NewDecoder(rr.Body).Decode(&keys); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if keys.PublicKey == "" || keys.PrivateKey == "" {
		t.Errorf("expected non-empty public and private keys, got %+v", keys)
	}
}
