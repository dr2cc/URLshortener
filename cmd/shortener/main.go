package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Storage interface {
	InsertURL(uid string, url string) error
	GetURL(uid string) (string, error)
}

// тип urlStorage и его параметр Data
type UrlStorage struct {
	Data map[string]string
}

// конструктор объектов с типом urlStorage
func NewStorageStruct() *UrlStorage {
	return &UrlStorage{
		Data: make(map[string]string),
	}
}

// тип urlStorage и его метод InsertURL
func (s *UrlStorage) InsertURL(uid string, url string) error {
	s.Data[uid] = url
	return nil
}

// тип urlStorage и его метод GetURL
func (s *UrlStorage) GetURL(uid string) (string, error) {
	e, existss := s.Data[uid]
	if !existss {
		return uid, errors.New("URL with such id doesn`t exist")
	}
	return e, nil
}

//*******************************************************************
// Реализую интерфейс Storage

func MakeEntry(s Storage, uid string, url string) {
	s.InsertURL(uid, url)
}

func GetEntry(s Storage, uid string) (string, error) {
	e, err := s.GetURL(uid)
	return e, err
}

//********************************************************************

func generateShortURL(urlList *UrlStorage, longURL string) string {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел
	runes := []rune(longURL)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	//удаляю из полученной строки все кроме букв и цифр
	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	//[:11] здесь сокращаю строку
	id := reg.ReplaceAllString(string(runes[:11]), "")

	//Реализую интерфейс Storage, что в последующем даст возможность
	//использовать его методы и другим типам
	MakeEntry(urlList, id, longURL)

	return "/" + id
}

// тип urlStorage и его метод PostHandler
func (ts *UrlStorage) PostHandler(w http.ResponseWriter, req *http.Request) {
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

			// Генерирую ответ и создаю запись в хранилище
			response := req.Host + generateShortURL(ts, string(param))

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, response)
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Content-Type isn`t text/plain")
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Method not allowed")
	}
}

// тип urlStorage и его метод GetHandler
func (ts *UrlStorage) GetHandler(w http.ResponseWriter, req *http.Request) {
	//Тесты подсказали добавить проверку на метод:
	switch req.Method {
	case http.MethodGet:
		// //Пока (14.04.2025) не знаю как передать PathValue при тестировании.
		// id := req.PathValue("id")

		// А вот RequestURI получается и от клиента и из теста
		// Но получаю лишний "/"
		id := strings.TrimPrefix(req.RequestURI, "/")

		// //Так не реализуя интерфейс
		//longURL, err := ts.GetURL(id)

		//Реализую интерфейс
		longURL, err := GetEntry(ts, id)

		if err != nil {
			//http.Error(w, "URL not found", http.StatusBadRequest)
			w.Header().Set("Location", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", longURL)
		// //И так и так работает. Оставил первоначальный вариант.
		//http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.Header().Set("Location", "Method not allowed")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// ******************************************************************************
// Секция переопредения стандартного ServeMux маршрутизатором CustomMux,
// Цель- возвращать 400 вместо 405

type CustomMux struct {
	*http.ServeMux
}

func (m *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Проверяю наличие пути
	_, pattern := m.Handler(r)
	if pattern == "" {
		// // Если эндпоинта нет вообще то стандарт- 404
		// http.NotFound(w, r)

		// Но мне нужно 400
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Если эндпоинт есть, но метод не совпадает- 400
	if !m.isMethodAllowed(r) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Иначе передаю обработку стандартному ServeMux
	m.ServeMux.ServeHTTP(w, r)
}

// isMethodAllowed проверяет, разрешен ли метод для данного пути
func (m *CustomMux) isMethodAllowed(r *http.Request) bool {
	// Получаю зарегистрированный обработчик для этого пути
	handler, _ := m.Handler(r)

	// Если обработчик — это ServeMux (значит, метод не совпадает)
	_, isServeMux := handler.(*http.ServeMux)
	return !isServeMux
}

//*************************************************************************

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
