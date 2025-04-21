package main

import (
	"github.com/dr2cc/URLshortener.git/internal/handlers"
	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	storageInstance := storage.NewStorage()

	router := gin.Default()

	router.POST("/", handlers.PostHandler(storageInstance))
	router.GET("/:id", handlers.GetHandler(storageInstance))

	router.Run("localhost:8080")
}
