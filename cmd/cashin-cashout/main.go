package main

import (
	"log"
	"net/http"
	"os"

	"cash-flow/src/application"
	"cash-flow/src/domain/transaction"
	"cash-flow/src/infrastructure/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	//err := godotenv.Load("cmd/cashin-cashout/.env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("cashin-cashout: Error loading .env file")
	}

	useAuth := os.Getenv("USE_KEYCLOAK")
	r := chi.NewRouter()
	db := database.NewDB()

	transactionUserCase := transaction.TransactionUseCase{
		Repository: &transaction.TransactionRepository{Db: db},
	}

	handler := application.HandlerTransaction{
		TransactionUseCase: &transactionUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", application.HealthCheck)
	r.Route("/transactions", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(application.Auth)
		}
		r.Post("/create", handler.CreateTransaction)
		r.Get("/", handler.GetTransactions)
	})

	http.ListenAndServe(":8088", r)
}
