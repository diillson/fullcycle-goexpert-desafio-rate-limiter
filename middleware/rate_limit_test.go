package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockLimiter struct {
	allowed bool
	err     error
}

func (m *MockLimiter) IsAllowed(ip string, token string) (bool, error) {
	return m.allowed, m.err
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		allowed        bool
		expectedStatus int
	}{
		{
			name:           "Allow request",
			allowed:        true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Block request",
			allowed:        false,
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLimiter := &MockLimiter{allowed: tt.allowed}
			middleware := NewRateLimitMiddleware(mockLimiter)

			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/api", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
