package main

import (
	"cash-flow/src/application"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rotas
	r.Get("/daily-summary", application.GetDailySummary)

	http.ListenAndServe(":8081", r)
}
