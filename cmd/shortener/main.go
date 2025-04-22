package main

import (
	"github.com/dr2cc/URLshortener.git/internal/handlers"
	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	storageInstance := storage.NewStorage()
	postHandler := handlers.PostHandler(storageInstance)
	getHandler := handlers.GetHandler(storageInstance)

	router := gin.Default()

	// Добавляем middleware проверки методов
	// Сервер должен возвращать только 400, на все не корректные запросы
	router.Use(handlers.MethodChecker())

	// Явно регистрируем только разрешенные методы
	router.Handle("POST", "/", postHandler)
	router.Handle("GET", "/:id", getHandler)

	router.Run("localhost:8080")
}
