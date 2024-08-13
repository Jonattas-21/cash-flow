package main

import (
	"cash-flow/src/application"
	"cash-flow/src/domain/dailySummary"
	"cash-flow/src/infrastructure/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cash-flow/cmd/daily-summary/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	db := database.NewDB()

	dailySummaryUserCase := dailySummary.DailySummaryUseCase{
		Repository:       &dailySummary.DailySummaryRepository{Db: db},
		CashinCashoutUrl: os.Getenv("CASHINCASHOUT"),
	}

	handler := application.HandlerSummary{
		DailySummaryUseCase: &dailySummaryUserCase,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", application.HealthCheck)
	r.Route("/transactions", func(r chi.Router) {
		r.Use(application.Auth)
		r.Get("/generateDailyReport", handler.GetDailySummary)
	})

	http.ListenAndServe(":8081", r)
}
