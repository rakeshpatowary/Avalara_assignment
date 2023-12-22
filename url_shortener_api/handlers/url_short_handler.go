package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var URLMap = make(map[string]string)

type URLShortenRequest struct {
	Destination string `json:"destination"`
}

type URLShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// URLShortener is the http handler to short any valid URL
func URLShortener(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		urlSortener(rw, r)
		return
	}

	// Catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// urlSortener generates short URL for valid input URL
func urlSortener(rw http.ResponseWriter, r *http.Request) {
	var usReq URLShortenRequest

	err := json.NewDecoder(r.Body).Decode(&usReq)
	if err != nil {
		http.Error(rw, "Invalid request format", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix((string)(usReq.Destination), "http") {
		http.Error(rw, "Not a valid URL ", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey(5)

	// Map the input URL to the generated short key
	URLMap[shortKey] = usReq.Destination

	shortURL := fmt.Sprintf("http://avalara-domain.com/%s", shortKey)

	response := URLShortenResponse{
		ShortURL: shortURL,
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

// RedirectToOriginalURL redirects to original URl based on the shortKey
func RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	varMap := mux.Vars(r)

	fmt.Println("varMap: ", varMap)
	shortKey := varMap["shortKey"]

	originalURL, exists := URLMap[shortKey]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// generateShortKey returns a slice of string of length keyLength
func generateShortKey(keyLength int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}
