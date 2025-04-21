package main

import (
	"net/http"
)

func main() {
	// mux := http.NewServeMux()

	//Для создания ответа 400 на все не верные запросы
	//создаю кастомный ServeMux (маршрутизатор)
	mux := &CustomMux{http.NewServeMux()}

	//создаю объект типа UrlStorage
	storage := NewStorageStruct()

	//обращаюсь к методам UrlStorage
	mux.HandleFunc("POST /{$}", storage.PostHandler)
	mux.HandleFunc("GET /{id}", storage.GetHandler)

	http.ListenAndServe("localhost:8080", mux)
}
