package main

import (
	"log"
	"net/http"

	"cash-flow/src/application"
	"cash-flow/src/domain/transaction"
	"cash-flow/src/infrastructure/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	db := database.NewDB()

	transactionUserCase := transaction.TransactionUseCase{
		IRepository: &database.Repository{Db: db},
	}

	handler := application.Handler{
		TransactionUseCase: &transactionUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", application.HealthCheck)
	r.Route("/transactions", func(r chi.Router) {
		r.Use(application.Auth)
		r.Post("/create", handler.CreateTransaction)
		r.Get("/", handler.GetTransactions)
	})

	http.ListenAndServe(":8080", r)
}
