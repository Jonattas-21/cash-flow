package usecases_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Jonattas-21/cash-flow/internal/usecases"
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"time"
	"github.com/stretchr/testify/mock"
	internalMock "github.com/Jonattas-21/cash-flow/tests"
)

var(
	repositoryMock = new(internalMock.TransactionRepositoryMock)
	userCase = usecases.TransactionUseCase{
		Repository: repositoryMock,
	}
	newTransactionDto = &entities.Transaction{
		Amount: 100,
		Type: "credit",
		Date: time.Now().Truncate(24*time.Hour),
	}
)

func Setup() {
	repositoryMock = new(internalMock.TransactionRepositoryMock)
	userCase = usecases.TransactionUseCase{
		Repository: repositoryMock,
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

