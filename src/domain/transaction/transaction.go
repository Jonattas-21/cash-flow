package transaction

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	Amount    int   `json:"amount"`
	Type      string    `json:"type"` // "credit" ou "debit"
	CreatedAt time.Time `json:"CreatedAt"`
	CreatedBy string    `json:"CreatedBy"`
}
