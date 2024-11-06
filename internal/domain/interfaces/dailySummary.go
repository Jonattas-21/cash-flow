package interfaces

import (
	"internal/domain/entities"
	"time"
)

type DailySummaryRepository interface {
	SaveReport(item entities.DailySummary) error
	GetReport(date time.Time) (*entities.DailySummary, error)
}
