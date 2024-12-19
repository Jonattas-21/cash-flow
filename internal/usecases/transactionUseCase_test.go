package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/usecases"
	internalMock "github.com/Jonattas-21/cash-flow/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	repositoryMock = new(internalMock.TransactionRepositoryMock)
	userCase       = usecases.TransactionUseCase{
		Repository: repositoryMock,
	}
	newTransactionDto = &entities.Transaction{}
)

func Setup() {
	repositoryMock = new(internalMock.TransactionRepositoryMock)
	userCase = usecases.TransactionUseCase{
		Repository: repositoryMock,
	}
	newTransactionDto = &entities.Transaction{
		Amount: 100,
		Type:   "credit",
		Date:   time.Now().Truncate(24 * time.Hour),
	}
}

func Test_SaveTransaction(t *testing.T) {
	Setup()
	assert := assert.New(t)

	repositoryMock.On("Save", mock.Anything).Return(nil)

	obj, validations, err := userCase.SaveTransaction(*newTransactionDto)

	assert.NotNil(obj)
	assert.Empty(validations)
	assert.Nil(err)
}

func Test_SaveTransactions_ValidationsError_amount(t *testing.T) {
	Setup()
	assert := assert.New(t)
	newTransactionDto.Amount = 0
	_, validations, err := userCase.SaveTransaction(*newTransactionDto)

	assert.Equal("Amount is required and must be greater than 0", validations[0])
	assert.Nil(err)
}

func Test_SaveTransactions_ValidationsError_type(t *testing.T) {
	Setup()
	assert := assert.New(t)
	newTransactionDto.Type = "dafault"
	_, validations, err := userCase.SaveTransaction(*newTransactionDto)

	assert.Equal("Transaction type is required and must be credit or debit", validations[0])
	assert.Nil(err)
}

func Test_SaveTransactions_ValidationsError_date(t *testing.T) {
	Setup()
	assert := assert.New(t)
	newTransactionDto.Date = time.Time{}
	_, validations, err := userCase.SaveTransaction(*newTransactionDto)

	assert.Equal("Date is required and must be greater than today", validations[0])
	assert.Nil(err)
}

func Test_FindAllTransactions(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("FindAll").Return([]entities.Transaction{}, nil)
	transactions, err := userCase.FindAllTransactions()

	assert.NotNil(transactions)
	assert.Nil(err)
}

func Test_FindAllTransactions_error(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("FindAll").Return([]entities.Transaction{}, errors.New("error"))
	_, err := userCase.FindAllTransactions()

	assert.NotNil(err)
}

func Test_FindTransactions(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("FindByDay").Return([]entities.Transaction{}, nil)
	date := time.Now().Truncate(24 * time.Hour)
	transactions, err := userCase.FindTransactions(date)

	assert.NotNil(transactions)
	assert.Nil(err)
}

func Test_FindTransactions_error(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("FindByDay").Return([]entities.Transaction{}, errors.New("error"))
	_, err := userCase.FindTransactions(time.Now().Truncate(24 * time.Hour))

	assert.NotNil(err)
}

func Test_DeleteTransaction(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("DeleteTransaction").Return(nil)
	err := userCase.DeleteTransaction("abc1")

	assert.Nil(err)
}

func Test_DeleteTransaction_error(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("DeleteTransaction").Return(errors.New("error"))
	err := userCase.DeleteTransaction("abc1")

	assert.NotNil(err)
}

func Test_UpdateTransaction(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("UpdateTransaction").Return(nil)
	err := userCase.UpdateTransaction("abc1", *newTransactionDto)

	assert.Nil(err)
}

func Test_UpdateTransaction_error(t *testing.T) {
	Setup()
	assert := assert.New(t)
	repositoryMock.On("UpdateTransaction").Return(errors.New("error"))
	err := userCase.UpdateTransaction("abc1", *newTransactionDto)

	assert.NotNil(err)
}

