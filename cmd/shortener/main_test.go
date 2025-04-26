package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLshortener.git/internal/handlers"
	"github.com/dr2cc/URLshortener.git/internal/server"
	"github.com/dr2cc/URLshortener.git/internal/storage"
)

// Integration test
func Test_main(t *testing.T) {
	// Подготовка (Setup)
	// * Создается реальное хранилище (URLStorage), а не mock.
	// * Инициализируются обработчики POST и GET, которые используют это хранилище.
	// * Создается сервер с этими обработчиками.
	storageInstance := storage.NewStorage()
	postHandler := handlers.PostHandler(storageInstance)
	getHandler := handlers.GetHandler(storageInstance)
	server := server.NewServer(postHandler, getHandler)

	t.Run("POST then GET", func(t *testing.T) {
		// Тест-кейс "POST then GET"
		// Шаг 1: POST-запрос для сокращения URL
		// Что проверяем:
		// * Обработчик PostHandler должен принять URL https://example.com, сохранить его в хранилище и вернуть сокращенную ссылку.
		// Как:
		// * Отправляется POST-запрос с телом text/plain.
		// * Сервер обрабатывает запрос через PostHandler.
		// Ожидаемый результат:
		// * HTTP-статус 201 Created.
		// * В теле ответа — сокращенный URL (например, http://example/abc123).
		postReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		postReq.Header.Set("Content-Type", "text/plain")
		postW := httptest.NewRecorder()
		server.ServeHTTP(postW, postReq)

		// Проверка (Assertion) POST- Проверяется, что сервер вернул статус 201 Created
		if postW.Result().StatusCode != http.StatusCreated {
			t.Fatalf("POST failed: %d", postW.Result().StatusCode)
		}

		// Шаг 2: Извлечение shortURL
		// Извлекаем только ID (всё после последнего '/'), т.е. из ответа (http://example/abc123) извлекается идентификатор (abc123).
		shortURL := strings.TrimSpace(postW.Body.String())
		id := shortURL[strings.LastIndex(shortURL, "/")+1:]
		t.Logf("Correct extracted ID: %s", id) // Должно быть "emtsapce"

		// Проверяем хранилище
		_, err := storageInstance.GetURL(id)
		if err != nil {
			t.Fatalf("URL not found in storage for ID '%s'. Storage contents: %+v",
				id, storageInstance.Data)
		}

		// Шаг 3: GET-запрос для получения исходного URL
		// Что проверяем:
		// * Обработчик GetHandler должен найти исходный URL по идентификатору abc123 и вернуть 307.
		// Как:
		// * Отправляется GET-запрос по пути /{id} (например, /abc123).
		// * Сервер обрабатывает запрос через GetHandler.
		// Ожидаемый результат:
		// * HTTP-статус 307 Temporary Redirect.
		// * В заголовке Location — исходный URL (https://example.com).
		getReq := httptest.NewRequest(http.MethodGet, "/"+id, nil)
		getW := httptest.NewRecorder()
		server.ServeHTTP(getW, getReq)

		// Проверка (Assertion) GET- проверяется статус 307 Temporary Redirect
		if getW.Result().StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("Expected status 307, got %d", getW.Result().StatusCode)
		}

		// Проверка (Assertion) GET- проверяется, что заголовок Location содержит исходный URL
		location := getW.Result().Header.Get("Location")
		if location != "https://example.com" {
			t.Errorf("Expected https://example.com, got %s", location)
		}
	})
}
