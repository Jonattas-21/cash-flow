package tests

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/stretchr/testify/mock"
	"time"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (r *TransactionRepositoryMock) Save(transaction entities.Transaction) error {
	args := r.Called(transaction)
	return args.Error(0)
}

func (r *TransactionRepositoryMock) FindAll() ([]entities.Transaction, error) {
	args := r.Called()
	return args.Get(0).([]entities.Transaction), args.Error(1)
}

func (r *TransactionRepositoryMock) FindByDay(date time.Time) ([]entities.Transaction, error){
	args := r.Called()
	return args.Get(0).([]entities.Transaction), args.Error(1)
}

func (r *TransactionRepositoryMock) DeleteTransaction(id string) error {
	args := r.Called()
	return args.Error(0)
}

func (r *TransactionRepositoryMock) UpdateTransaction(id string, transaction entities.Transaction) error {
	args := r.Called()
	return args.Error(0)
}