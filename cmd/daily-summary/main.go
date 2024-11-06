package main

import (
	"internal/api"
	"internal/api/handlers"
	"internal/domain/entities"
	"internal/infrastructure/database"
	"internal/infrastructure/repositories"
	"internal/usecases"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cash-flow/cmd/daily-summary/.env")
	useAuth := os.Getenv("USE_KEYCLOAK")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	db := database.NewDB()

	// Migrate the schema
	db.AutoMigrate(&entities.DailySummary{})

	rdb := database.NewCache()

	dailySummaryUserCase := usecases.DailySummaryUseCase{
		Repository:       &repositories.DailySummaryRepository{Db: db},
		CashinCashoutUrl: os.Getenv("CASHINCASHOUT"),
		Rdb:              rdb,
	}

	handler := handlers.HandlerSummary{
		DailySummaryUseCase: &dailySummaryUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.HealthCheck)
	r.Route("/transactions", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(api.Auth)
		}
		r.Get("/generateDailyReport", handler.GetDailySummary)
	})

	http.ListenAndServe(":8081", r)
}
