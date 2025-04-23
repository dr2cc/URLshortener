package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dr2cc/URLshortener.git/internal/storage"
)

func generateShortURL(urlList *storage.UrlStorage, longURL string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := []rune(longURL)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	//[:11] здесь сокращаю строку
	id := reg.ReplaceAllString(string(runes[:11]), "")

	storage.MakeEntry(urlList, id, longURL)

	return "/" + id
}

func PostHandler(ts *storage.UrlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			switch req.Header.Get("Content-Type") {
			case "text/plain":
				//param - тело запроса (тип []byte)
				param, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				//longURL := string(param)
				response := req.Host + generateShortURL(ts, string(param))

				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, response)
			default:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Content-Type isn't text/plain")
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Method not allowed")
		}
	}
}

func GetHandler(ts *storage.UrlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			id := strings.TrimPrefix(req.RequestURI, "/")
			longURL, err := storage.GetEntry(ts, id)
			if err != nil {
				w.Header().Set("Location", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Location", longURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.Header().Set("Location", "Method not allowed")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
