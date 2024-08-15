package transaction

import "time"

type Transaction struct {
	ID        string    `json:"id" gorm:"size:50"`
	Amount    int       `json:"amount"`
	Type      string    `json:"type"` // "credit" ou "debit"
	Date	  time.Time `json:"date"`
	CreatedAt time.Time `json:"CreatedAt"`
	CreatedBy string    `json:"CreatedBy" gorm:"size:50"`
}
