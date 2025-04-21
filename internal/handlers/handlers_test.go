package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(ts *storage.UrlStorage) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", PostHandler(ts))
	router.GET("/:id", GetHandler(ts))
	return router
}

func TestPostHandler(t *testing.T) {
	store := storage.NewStorage()
	router := setupRouter(store)

	t.Run("Successful URL shortening", func(t *testing.T) {
		body := bytes.NewBufferString("https://example.com")
		req, _ := http.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "/")
	})

	t.Run("Invalid content type", func(t *testing.T) {
		body := bytes.NewBufferString("https://example.com")
		req, _ := http.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type isn't text/plain")
	})

	t.Run("Wrong HTTP method", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetHandler(t *testing.T) {
	store := storage.NewStorage()
	store.InsertURL("abc123", "https://example.com")
	router := setupRouter(store)

	t.Run("Successful redirect", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/abc123", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Location"))
	})

	t.Run("Non-existent URL", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/nonexistent", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Wrong HTTP method", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/abc123", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGenerateShortURL(t *testing.T) {
	store := storage.NewStorage()
	longURL := "https://example.com/very/long/url"

	shortURL := generateShortURL(store, longURL)
	assert.NotEmpty(t, shortURL)
	assert.True(t, len(shortURL) > 1) // At least "/" + one character

	// Verify the URL was stored
	id := shortURL[1:] // Remove leading "/"
	url, err := store.GetURL(id)
	assert.NoError(t, err)
	assert.Equal(t, longURL, url)
}
