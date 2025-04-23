package main

import (
	"net/http"

	"github.com/dr2cc/URLshortener.git/internal/handlers"
	"github.com/dr2cc/URLshortener.git/internal/server"
	"github.com/dr2cc/URLshortener.git/internal/storage"
)

func main() {
	storageInstance := storage.NewStorage()

	postHandler := handlers.PostHandler(storageInstance)
	getHandler := handlers.GetHandler(storageInstance)

	server := server.NewServer(postHandler, getHandler)
	http.ListenAndServe(":8080", server)
}
