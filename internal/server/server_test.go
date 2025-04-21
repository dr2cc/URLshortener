package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomMux_ServeHTTP(t *testing.T) {
	mux := &CustomMux{http.NewServeMux()}
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "Valid route",
			method:     http.MethodGet,
			path:       "/test",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid method",
			method:     http.MethodPost,
			path:       "/test",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent route",
			method:     http.MethodGet,
			path:       "/nonexistent",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	postHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := NewServer(postHandler, getHandler)

	// Test POST route
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Error("POST route not working")
	}

	// Test GET route
	req = httptest.NewRequest(http.MethodGet, "/test123", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Error("GET route not working")
	}
}
