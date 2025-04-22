package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dr2cc/URLshortener.git/internal/handlers"
	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShortenerIntegration(t *testing.T) {
	// 1. ИНИЦИАЛИЗАЦИЯ ТЕСТОВОГО СЕРВЕРА
	// ----------------------------------
	// Устанавливаем тестовый режим Gin для отключения лишнего логгирования
	gin.SetMode(gin.TestMode)

	// Создаем тестовый роутер с обработчиками как в production
	router := setupTestRouter()

	// 2. ТЕСТ СОЗДАНИЯ КОРОТКОЙ ССЫЛКИ
	// ---------------------------------
	t.Run("Create short URL", func(t *testing.T) {
		// 2.1 ПОДГОТОВКА ЗАПРОСА
		// Создаем тело запроса с исходным URL
		body := bytes.NewBufferString("https://example.com")

		// Формируем POST-запрос к корневому эндпоинту
		req, _ := http.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", "text/plain") // Устанавливаем нужный заголовок

		// 2.2 ВЫПОЛНЕНИЕ ЗАПРОСА
		// Создаем Recorder для записи ответа
		w := httptest.NewRecorder()

		// Передаем запрос в роутер
		router.ServeHTTP(w, req)

		// 2.3 ПРОВЕРКА РЕЗУЛЬТАТОВ
		// Проверяем код статуса (должен быть 201 Created)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Проверяем, что в ответе есть короткий URL (должен содержать "/")
		assert.Contains(t, w.Body.String(), "/")
	})

	// 3. ТЕСТ РЕДИРЕКТА ПО КОРОТКОЙ ССЫЛКЕ
	// -------------------------------------
	t.Run("Redirect by short URL", func(t *testing.T) {
		// 3.1 ПОДГОТОВКА ТЕСТОВЫХ ДАННЫХ
		// Исходный URL для сокращения
		testURL := "https://google.com"

		// 3.2 СОЗДАНИЕ КОРОТКОЙ ССЫЛКИ
		// Формируем POST-запрос как в предыдущем тесте
		body := bytes.NewBufferString(testURL)
		reqPost, _ := http.NewRequest("POST", "/", body)
		reqPost.Header.Set("Content-Type", "text/plain")

		// Выполняем запрос
		wPost := httptest.NewRecorder()
		router.ServeHTTP(wPost, reqPost)

		// 3.3 ИЗВЛЕЧЕНИЕ ID КОРОТКОЙ ССЫЛКИ
		// Полный ответ сервера (формат "localhost:8080/abc123")
		fullShortURL := wPost.Body.String()

		// Разбиваем URL по символу "/" чтобы извлечь ID
		parts := strings.Split(fullShortURL, "/")

		// Берем последний элемент (собственно ID)
		id := parts[len(parts)-1]

		// Проверяем что ID не пустой
		assert.NotEmpty(t, id, "ID не должен быть пустым")

		// 3.4 ПРОВЕРКА РЕДИРЕКТА
		// Формируем GET-запрос к эндпоинту /{id}
		reqGet, _ := http.NewRequest("GET", "/"+id, nil)

		// Создаем новый Recorder для этого запроса
		wGet := httptest.NewRecorder()

		// Выполняем запрос
		router.ServeHTTP(wGet, reqGet)

		// 3.5 ВАЛИДАЦИЯ РЕЗУЛЬТАТОВ
		// Проверяем код статуса (должен быть 307 Temporary Redirect)
		assert.Equal(t, http.StatusTemporaryRedirect, wGet.Code)

		// Проверяем заголовок Location (должен содержать исходный URL)
		assert.Equal(t, testURL, wGet.Header().Get("Location"))
	})
}

// setupTestRouter создает тестовый экземпляр роутера
// ------------------------------------------------
func setupTestRouter() *gin.Engine {
	// 1. Инициализация хранилища (как в production)
	storageInstance := storage.NewStorage()

	// 2. Создание обработчиков
	postHandler := handlers.PostHandler(storageInstance)
	getHandler := handlers.GetHandler(storageInstance)

	// 3. Настройка роутера Gin
	router := gin.Default()

	// Регистрируем эндпоинты:
	// POST / - создание короткой ссылки
	router.POST("/", postHandler)

	// GET /{id} - редирект по короткой ссылке
	router.GET("/:id", getHandler)

	return router
}
