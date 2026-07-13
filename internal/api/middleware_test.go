package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIKeyRateLimitAllowed(t *testing.T) {
	s := NewServer(":0", &AuthConfig{Enabled: false})
	s.RegisterAPIKey("test-key-123", 10)

	handler := s.handlerWithMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.Header.Set("X-API-Key", "test-key-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAPIKeyRateLimitExceeded(t *testing.T) {
	s := NewServer(":0", &AuthConfig{Enabled: false})
	s.RegisterAPIKey("limited-key", 2)

	handler := s.handlerWithMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/v1/health", nil)
		req.Header.Set("X-API-Key", "limited-key")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if i < 2 && w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
		if i == 2 && w.Code != http.StatusTooManyRequests {
			t.Errorf("request 2: expected 429, got %d", w.Code)
		}
	}
}

func TestFallbackToIPBasedLimiter(t *testing.T) {
	s := NewServer(":0", &AuthConfig{Enabled: false})

	handler := s.handlerWithMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUnknownAPIKeyFallsBackToIP(t *testing.T) {
	s := NewServer(":0", &AuthConfig{Enabled: false})
	s.RegisterAPIKey("known-key", 5)

	handler := s.handlerWithMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.Header.Set("X-API-Key", "unknown-key-value")
	req.RemoteAddr = "10.0.0.1:9999"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for unknown API key fallback, got %d", w.Code)
	}

	_ = time.Now()
}
