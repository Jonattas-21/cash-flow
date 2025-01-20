package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/infrastructure/cache"
	"github.com/Jonattas-21/cash-flow/internal/usecases"
	internalMock "github.com/Jonattas-21/cash-flow/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	DailyRepositoryMock = new(internalMock.DailySummaryRepositoryMock)
	DailyUserCase       = usecases.DailySummaryUseCase{
		Repository: DailyRepositoryMock,
	}
	newDailyReportDto = &entities.DailySummary{}
)

func DailyReportSetup() {
	repositoryMock = new(internalMock.TransactionRepositoryMock)
	userCase = usecases.TransactionUseCase{
		Repository: repositoryMock,
	}

	DailyRepositoryMock = new(internalMock.DailySummaryRepositoryMock)
	DailyUserCase = usecases.DailySummaryUseCase{
		Repository:         DailyRepositoryMock,
		TransactionUseCase: &userCase,
		Rdb:                cache.NewCache(),
	}
	newDailyReportDto = &entities.DailySummary{
		Date:      time.Now(),
		Credit:    100,
		Debit:     0,
		Total:     100,
		Status:    "closed",
		CreatedAt: time.Now(),
	}
}

func Test_GetDailySummary_Partial(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	DailyRepositoryMock.On("DeleteReport", mock.Anything).Return(nil)
	DailyRepositoryMock.On("GetReport", mock.MatchedBy(func(date time.Time) bool {
		today := time.Now()
		if date.Day() == today.Day() && date.Month() == today.Month() && date.Year() == today.Year() {
			newDailyReportDto.Status = "partial"
		}

		return true
	})).Return(newDailyReportDto, nil)

	obj, err := DailyUserCase.GetDailySummary(time.Now())

	t.Log("obj", obj)
	assert.Equal(newDailyReportDto, obj)
	assert.Equal(obj.Status, "partial")
	assert.Nil(err)
}

func Test_GetDailySummary_error(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	DailyRepositoryMock.On("DeleteReport", mock.Anything).Return(nil)
	DailyRepositoryMock.On("GetReport", mock.Anything).Return(newDailyReportDto, errors.New("error"))

	obj, err := DailyUserCase.GetDailySummary(time.Now())

	assert.Nil(obj)
	assert.NotNil(err)
}

func Test_GetDailySummary_delete_report_error(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	DailyRepositoryMock.On("DeleteReport", mock.Anything).Return(errors.New("error delete report"))
	obj, err := DailyUserCase.GetDailySummary(time.Now())

	assert.Nil(obj)
	assert.NotNil(err)
}

func Test_GetDailySummary_Generate_report_error(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	transaction := append([]entities.Transaction{}, entities.Transaction{
		ID:        "1",
		Amount:    100,
		Type:      "credit",
		Date:      time.Now(),
		CreatedAt: time.Now(),
	})

	DailyRepositoryMock.On("DeleteReport", mock.Anything).Return(nil)
	DailyRepositoryMock.On("GetReport", mock.Anything).Return(&entities.DailySummary{Total: 0}, nil)
	repositoryMock.On("FindByDay", mock.Anything).Return(transaction, errors.New("error geting transactions"))

	_, err := DailyUserCase.GetDailySummary(time.Now())
	assert.NotNil(err)
}

func Test_GenerateReport(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	transaction := append([]entities.Transaction{}, entities.Transaction{
		ID:        "1",
		Amount:    100,
		Type:      "credit",
		Date:      time.Now(),
		CreatedAt: time.Now(),
	})

	repositoryMock.On("FindByDay", mock.Anything).Return(transaction, nil)
	DailyRepositoryMock.On("SaveReport", mock.Anything).Return(nil)

	obj, err := DailyUserCase.GenerateReport(time.Now())

	assert.Nil(err)
	assert.NotNil(obj)
}

func Test_GenerateReport_error(t *testing.T) {
	DailyReportSetup()
	assert := assert.New(t)

	repositoryMock.On("FindByDay", mock.Anything).Return([]entities.Transaction{}, nil)
	DailyRepositoryMock.On("SaveReport", mock.Anything).Return(errors.New("error"))

	obj, err := DailyUserCase.GenerateReport(time.Now())

	assert.NotNil(err)
	assert.Nil(obj)
}
