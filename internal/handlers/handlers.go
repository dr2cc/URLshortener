package handlers

import (
	"io"
	"math/rand"
	"regexp"
	"time"

	"github.com/dr2cc/URLshortener.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func generateShortURL(urlList *storage.UrlStorage, longURL string) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune(longURL)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	reg := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`)
	id := reg.ReplaceAllString(string(runes[:11]), "")

	storage.MakeEntry(urlList, id, longURL)

	return "/" + id
}

func PostHandler(ts *storage.UrlStorage) gin.HandlerFunc {
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

func GetHandler(ts *storage.UrlStorage) gin.HandlerFunc {
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
