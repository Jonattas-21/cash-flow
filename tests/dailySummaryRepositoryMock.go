package tests

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/stretchr/testify/mock"
	"time"
	)	 

type DailySummaryRepositoryMock struct {
	mock.Mock
}

func (d* DailySummaryRepositoryMock) SaveReport(item entities.DailySummary) error {
	args := d.Called(item)
	return args.Error(0)
}

func (d* DailySummaryRepositoryMock) GetReport(date time.Time) (entities.DailySummary, error) {
	args := d.Called(date)
	return args.Get(0).(entities.DailySummary), args.Error(1)
}

func (d* DailySummaryRepositoryMock) DeleteReport(date time.Time) error {
	args := d.Called(date)
	return args.Error(0)
}