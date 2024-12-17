package interfaces

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"time"
)

type TransactionRepository interface {
	Save(item entities.Transaction) error
	FindByDay(date time.Time) ([]entities.Transaction, error)
	FindAll() ([]entities.Transaction, error)
	DeleteTransaction(id string) error
	UpdateTransaction(id string, item entities.Transaction) error
}
