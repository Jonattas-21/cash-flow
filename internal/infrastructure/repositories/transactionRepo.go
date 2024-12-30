package repositories

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	Db *gorm.DB
}

func (r *TransactionRepository) Save(item entities.Transaction) error {
	return r.Db.Create(item).Error
}

func (r *TransactionRepository) FindByDay(date time.Time) ([]entities.Transaction, error) {
	var itens []entities.Transaction

	dateFormmatedmin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dateFormmatedMax := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	tx := r.Db.Where("date BETWEEN ? and ?", dateFormmatedmin, dateFormmatedMax).Find(&itens)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return itens, nil
}

func (r *TransactionRepository) FindAll() ([]entities.Transaction, error) {
	var itens []entities.Transaction
	tx := r.Db.Find(&itens)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return itens, nil
}

func (r *TransactionRepository) DeleteTransaction(id string) error {
	tx := r.Db.Where("id = ?", id).Delete(&entities.Transaction{})
	if tx.RowsAffected == 0 {
		return errors.New("TransactionID not found")
	}
	return tx.Error
}

func (r *TransactionRepository) UpdateTransaction(id string, item entities.Transaction) error {
	tx := r.Db.Where("id = ?", id).Updates(item)
	if tx.RowsAffected == 0 {
		return errors.New("TransactionID not found")
	}
	return tx.Error
}