package application

import (
	"cash-flow/src/domain/dailySummary"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HandlerSummary struct {
	DailySummaryUseCase dailySummary.IDailySummaryUseCase
}

// GetDailySummary retorna o consolidado diário
func (h *HandlerSummary) GetDailySummary(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		http.Error(w, "Missing 'date' query parameter", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	report, err := h.DailySummaryUseCase.GenerateReport(date)
	if err != nil {
		log.Fatalln("Error to get report: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}
