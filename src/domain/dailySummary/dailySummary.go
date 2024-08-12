package dailySummary

import "time"

type DailySummary struct {
	Date   time.Time `json:"date"`
	Credit float64   `json:"credit"`
	Debit  float64   `json:"debit"`
	Total  float64   `json:"total"`
}
