package dailySummary

import (
	"time"

	"gorm.io/gorm"
)

type IDailySummaryRepository interface {
	SaveReport(item DailySummary) error
	GetReport(date time.Time) (DailySummary, error)
}

type DailySummaryRepository struct {
	Db *gorm.DB
}

func (d *DailySummaryRepository) SaveReport(item DailySummary) error {
	return d.Db.Create(item).Error
}

func (d *DailySummaryRepository) GetReport(date time.Time) (DailySummary, error) {
	//todo
	return DailySummary{}, nil
}
