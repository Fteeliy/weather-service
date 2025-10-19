package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
