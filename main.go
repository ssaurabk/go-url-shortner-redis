package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Change this for cloud deployment
})

func generateShortCode() string {
	b := make([]byte, 6) // 6-byte random string
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:6] // Shorten to 6 characters
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode()
	err := rdb.Set(ctx, shortCode, originalURL, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Short URL: http://localhost:8083/%s", shortCode) // Change domain for production
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := rdb.Get(ctx, shortCode).Result()

	if err == redis.Nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenURL)
	http.HandleFunc("/", redirectURL)

	fmt.Println("Server is running on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
