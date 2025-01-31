package web_test

import (
	"net/http"
	"net/http/httptest"
	"sender/internal/server/web"
	"testing"
)

func TestWebServerRoutes(t *testing.T) {
	server := web.New("8080")

	// Тестируем маршруты
	tests := []struct {
		method         string
		url            string
		expectedStatus int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/keys/generate", http.StatusOK},
		{"POST", "/keys/generate", http.StatusMethodNotAllowed},
		{"GET", "/nonexistent", http.StatusNotFound},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.url, nil)
		rr := httptest.NewRecorder()

		// Используем server.routes вместо несуществующего метода Router()
		server.Routes().ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("for %s %s, expected status %d, got %d", test.method, test.url, test.expectedStatus, rr.Code)
		}
	}
}
