package interfaces

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"time"
)

type DailySummaryRepository interface {
	SaveReport(item entities.DailySummary) error
	GetReport(date time.Time) (entities.DailySummary, error)
	DeleteReport(date time.Time) error
}
