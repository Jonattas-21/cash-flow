package transaction

import (
	"errors"
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

	dateFormmatedmin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	dateFormmatedMax := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)

	tx := r.Db.Where("date BETWEEN ? and ?", dateFormmatedmin, dateFormmatedMax).Find(&itens)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return itens, nil
}
