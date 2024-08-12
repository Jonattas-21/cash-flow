package transaction

import (
	"time"

	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Save(item Transaction) error
	FindByDay(date time.Time) ([]Transaction, error)
}

type TransactionRepository struct {
	Db *gorm.DB
}

func (r *TransactionRepository) Save(item Transaction) error {
	return r.Db.Create(item).Error
}

func (r *TransactionRepository) FindByDay(date time.Time) ([]Transaction, error) {
	var itens []Transaction
	//err := r.Db.Find(&itens).Error
	return itens, nil
}
