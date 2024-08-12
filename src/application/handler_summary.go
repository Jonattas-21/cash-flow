package application

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"cashflow/internal/models"
)

var (
	transactions []models.Transaction
	mu           sync.Mutex
)

// GetDailySummary retorna o consolidado diário
func GetDailySummary(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	summaryMap := make(map[string]*models.DailySummary)

	for _, t := range transactions {
		date := t.Timestamp.Format("2006-01-02")
		if summaryMap[date] == nil {
			summaryMap[date] = &models.DailySummary{
				Date: t.Timestamp,
			}
		}
		if t.Type == "credit" {
			summaryMap[date].Credit += t.Amount
		} else {
			summaryMap[date].Debit += t.Amount
		}
		summaryMap[date].Total = summaryMap[date].Credit - summaryMap[date].Debit
	}

	var summaries []models.DailySummary
	for _, summary := range summaryMap {
		summaries = append(summaries, *summary)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}
