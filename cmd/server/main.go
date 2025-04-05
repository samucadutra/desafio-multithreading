package main

import (
	"desafio-multithreading/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	cepHandler := handlers.NewCepHandler()

	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// routes
	r.Route("/cep", func(r chi.Router) {
		r.Get("/{cep}", cepHandler.GetAddress)
	})

	http.ListenAndServe(":8000", r)
}
