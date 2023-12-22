package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url_shortener_api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "url-shortner", log.LstdFlags)
	sm := mux.NewRouter()

	sm.HandleFunc("/shortURL", handlers.URLShortener).Methods("PUT")
	sm.HandleFunc("/{shortKey}", handlers.RedirectToOriginalURL).Methods("GET")

	http.Handle("/", sm)

	// Create the server
	s := http.Server{
		Addr:    "localhost:9090",
		Handler: sm,
	}

	go func() {
		l.Println("Server is up & running on port", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc) // Graceful shutdown
}
