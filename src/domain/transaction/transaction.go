package transaction

import "time"

// Transação representa uma única transação financeira
type Transaction struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // "credit" ou "debit"
	CreatedAt time.Time `json:"CreatedAt"`
	CreatedBy string    `json:"CreatedBy"`
}
