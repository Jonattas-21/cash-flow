package interfaces

import (
	"cash-flow/internal/domain/entities"
	"time"
)

type DailySummaryRepository interface {
	SaveReport(item entities.DailySummary) error
	GetReport(date time.Time) (*entities.DailySummary, error)
}
