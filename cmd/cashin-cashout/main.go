package main

import (
	"log"
	"net/http"
	"os"

	"internal/api/handlers"
	"internal/domain/entities"
	"internal/infrastructure/database"
	"internal/infrastructure/repositories"
	"internal/usecases"

	"internal/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cmd/cashin-cashout/.env")
	if err != nil {
		log.Fatal("cashin-cashout: Error loading .env file")
	}

	useAuth := os.Getenv("USE_KEYCLOAK")

	//Create a new router
	r := chi.NewRouter()

	//Conect and migrate the schema
	db := database.NewDB()
	db.AutoMigrate(&entities.Transaction{})

	transactionUserCase := usecases.TransactionUseCase{
		Repository: &repositories.TransactionRepository{Db: db},
	}

	handler := handlers.HandlerTransaction{
		TransactionUseCase: &transactionUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.HealthCheck)
	r.Route("/transactions", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(api.Auth)
		}
		r.Post("/create", handler.CreateTransaction)
		r.Get("/", handler.GetTransactions)
	})

	http.ListenAndServe(":8088", r)
}
