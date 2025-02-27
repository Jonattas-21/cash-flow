// Package classification of API.
//
// Documentation of our API.
//
//	Schemes: http, https
//	BasePath: /api/v1
//	Version: 1.0.0
//	Host: localhost:8088
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jonattas-21/cash-flow/internal/api"
	"github.com/Jonattas-21/cash-flow/internal/api/handlers"
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/infrastructure/cache"
	"github.com/Jonattas-21/cash-flow/internal/infrastructure/database"
	"github.com/Jonattas-21/cash-flow/internal/infrastructure/repositories"
	"github.com/Jonattas-21/cash-flow/internal/usecases"

	_ "github.com/Jonattas-21/cash-flow/cmd/cashin-cashout/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	err := godotenv.Load("cmd/cashin-cashout/.env")
	if err != nil {
		log.Fatal("cashin-cashout: Error loading .env file")
	}

	useAuth := os.Getenv("USE_KEYCLOAK")

	//Create a new router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//Setup cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//Conect and migrate the schema
	db := database.NewDB()
	db.AutoMigrate(&entities.Transaction{}, &entities.DailySummary{})

	transactionUserCase := usecases.TransactionUseCase{
		Repository: &repositories.TransactionRepository{Db: db},
	}

	//Conect to Redis cache
	rdb := cache.NewCache()

	dailySummaryUserCase := usecases.DailySummaryUseCase{
		Repository:       &repositories.DailySummaryRepository{Db: db},
		TransactionUseCase: &transactionUserCase,
		CashinCashoutUrl: "",
		Rdb:              rdb,
	}

	handler := handlers.HandlerTransaction{
		TransactionUseCase: &transactionUserCase,
	}

	handlerSummary := handlers.HandlerSummary{
		DailySummaryUseCase: &dailySummaryUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", handlers.HealthCheck)

	//Transaction routes
	r.Route("/transactions", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(api.Auth)
		}
		r.Post("/create", handler.CreateTransaction)
		r.Patch("/update/{id}", handler.UpdateTransaction)
		r.Delete("/delete/{id}", handler.DeleteTransaction)
		r.Get("/", handler.GetTransactions)
	})

	//Daily Summary routes
	r.Route("/report", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(api.Auth)
		}
		r.Get("/getDailyReport", handlerSummary.GetDailySummary)
	})

	// Serve the Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8088/swagger/doc.json"))) // The url pointing to API definition

	http.ListenAndServe(":8088", r)
}
