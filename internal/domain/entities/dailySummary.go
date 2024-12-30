package entities

import (
	"time"
)

type DailySummary struct {
	Date   time.Time `json:"date"`
	Credit int       `json:"credit"`
	Debit  int       `json:"debit"`
	Total  int       `json:"total"`
	Status string    `json:"status" gorm:"size:50"`
	CreatedAt time.Time `json:"CreatedAt"`
}
