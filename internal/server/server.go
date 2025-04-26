package server

import (
	"net/http"
)

type CustomMux struct {
	*http.ServeMux
}

func (m *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := m.Handler(r)
	if pattern == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !m.isMethodAllowed(r) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	m.ServeMux.ServeHTTP(w, r)
}

func (m *CustomMux) isMethodAllowed(r *http.Request) bool {
	handler, _ := m.Handler(r)
	_, isServeMux := handler.(*http.ServeMux)
	return !isServeMux
}

func NewServer(postHandler, getHandler http.HandlerFunc) *CustomMux {
	mux := &CustomMux{http.NewServeMux()}
	mux.HandleFunc("POST /", postHandler)
	mux.HandleFunc("GET /{id}", getHandler)
	return mux
}
