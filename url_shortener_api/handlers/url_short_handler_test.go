package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var SuccessfulShortenTests = []struct {
	reqMap    map[string]string
	wantEqual string
}{
	{
		reqMap:    map[string]string{"destination": "https://www.google.com/"},
		wantEqual: "http://avalara-domain.com/",
	},
	{
		reqMap:    map[string]string{"destination": "https://pkg.go.dev/"},
		wantEqual: "http://avalara-domain.com/",
	},
	{
		reqMap:    map[string]string{"destination": "https://go.dev/"},
		wantEqual: "http://avalara-domain.com/",
	},
}

func TestURLShortener(t *testing.T) {
	for _, test := range SuccessfulShortenTests {
		router := mux.NewRouter()

		router.HandleFunc("/shortURL", URLShortener).Methods(http.MethodPut)

		// Testcase: Successful scenario
		t.Run("URLShortener_Success", func(t *testing.T) {
			var reqBody = test.reqMap

			jsonBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("Test couldn't be proceed , JSON marshal error: %v", err)
				return
			}
			req := httptest.NewRequest(http.MethodPut, "/shortURL", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to record the response
			rec := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status as %v, but got %v", http.StatusOK, rec.Code)
			}

			// Test the response
			var response URLShortenResponse
			err = json.NewDecoder(rec.Body).Decode(&response)
			if err != nil {
				t.Errorf("Error decoding the JSON response: %v", err)
			}

			if !strings.HasPrefix(response.ShortURL, test.wantEqual) {
				t.Errorf("Expected short URL start with domain %v, but got %v", test.wantEqual, response.ShortURL)
			}
		})
	}

	// Not succesful scenario: [Invalid JSON format in request body]
	t.Run("Invalid JSON Format", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"invalid_key": "invalid_value"}`)
		req, err := http.NewRequest("PUT", "/shortURL", reqBody)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rec := httptest.NewRecorder()

		// Call the URLShortener function
		URLShortener(rec, req)

		// Check the response code
		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %v, but got %v", http.StatusBadRequest, rec.Code)
		}
	})
}

var RedirectOriginalURLTests = []struct {
	shortKey  string
	wantEqual string
}{
	{
		shortKey:  "abcde",
		wantEqual: "https://www.google.com/",
	},
	{
		shortKey:  "efghi",
		wantEqual: "https://pkg.go.dev/",
	},
	{
		shortKey:  "jklmn",
		wantEqual: "https://go.dev/",
	},
}

func TestRedirectToOriginalURL(t *testing.T) {

	for _, test := range RedirectOriginalURLTests {
		router := mux.NewRouter()
		router.HandleFunc("/{shortKey}", RedirectToOriginalURL).Methods(http.MethodGet)

		// Test case: Successful redirection
		t.Run("RedirectToOriginalURL_Success", func(t *testing.T) {
			URLMap[test.shortKey] = test.wantEqual

			// Create request with the short key
			request := httptest.NewRequest(http.MethodGet, "/"+test.shortKey, nil)

			// Create a response recorder to record the response
			rec := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rec, request)

			// Check the status code
			if rec.Code != http.StatusMovedPermanently {
				t.Errorf("Expected status code %v, but got %v", http.StatusMovedPermanently, rec.Code)
			}
		})
	}
}
