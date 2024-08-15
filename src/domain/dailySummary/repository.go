package dailySummary

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type IDailySummaryRepository interface {
	SaveReport(item DailySummary) error
	GetReport(date time.Time) (*DailySummary, error)
}

type DailySummaryRepository struct {
	Db *gorm.DB
}

func (d *DailySummaryRepository) SaveReport(item DailySummary) error {
	return d.Db.Create(item).Error
}

func (d *DailySummaryRepository) GetReport(date time.Time) (*DailySummary, error) {
	var summary DailySummary
	date = date.Truncate(24 * time.Hour)
	tx := d.Db.Where("date = ?", date).First(&summary)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &summary, tx.Error
}
