package handlers

import (
	"io"
	"math/rand"
	"regexp"
	"time"

	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func generateShortURL(urlList *storage.URLStorage, longURL string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := []rune(longURL)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	id := reg.ReplaceAllString(string(runes[:11]), "")

	storage.MakeEntry(urlList, id, longURL)

	return "/" + id
}

func PostHandler(ts *storage.URLStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Method not allowed"})
			return
		}

		if c.ContentType() != "text/plain" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Content-Type isn't text/plain"})
			return
		}

		param, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		response := c.Request.Host + generateShortURL(ts, string(param))
		c.String(201, response)
	}
}

func GetHandler(ts *storage.URLStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Method not allowed"})
			return
		}

		id := c.Param("id")
		longURL, err := storage.GetEntry(ts, id)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(307, longURL)
	}
}

// middleware проверки методов
func MethodChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем список разрешенных методов для текущего пути
		allowedMethods := getAllowedMethods(c)

		// Проверяем, разрешен ли текущий метод
		valid := false
		for _, method := range allowedMethods {
			if method == c.Request.Method {
				valid = true
				break
			}
		}

		if !valid {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Bad Request - invalid method",
			})
			return
		}

		c.Next()
	}
}

// getAllowedMethods возвращает разрешенные методы для текущего пути
func getAllowedMethods(c *gin.Context) []string {
	// В Gin нет прямого способа получить разрешенные методы,
	// поэтому мы используем обходной путь - проверяем путь вручную
	path := c.FullPath()

	switch path {
	case "/":
		return []string{"POST"}
	case "/:id":
		return []string{"GET"}
	default:
		return []string{}
	}
}
