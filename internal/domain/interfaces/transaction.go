package interfaces

import (
	"cash-flow/internal/domain/entities"
	"time"
)

type TransactionRepository interface {
	Save(item entities.Transaction) error
	FindByDay(date time.Time) ([]entities.Transaction, error)
}
