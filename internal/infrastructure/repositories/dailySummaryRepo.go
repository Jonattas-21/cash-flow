package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/Jonattas-21/cash-flow/internal/domain/entities"

	"gorm.io/gorm"
)

type DailySummaryRepository struct {
	Db *gorm.DB
}

func (d *DailySummaryRepository) GetReport(date time.Time) (entities.DailySummary, error) {
	var summary []entities.DailySummary

	dateFormmatedmin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dateFormmatedMax := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	tx := d.Db.Where("date BETWEEN ? and ?", dateFormmatedmin, dateFormmatedMax).Find(&summary)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || len(summary) == 0 {
		log.Println("DailySummaryRepository: GetReport: Error: ", tx.Error)
		return entities.DailySummary{}, nil
	}

	log.Println("DailySummaryRepository: GetReport: Record found")
	return summary[0], tx.Error
}

func (d *DailySummaryRepository) SaveReport(item entities.DailySummary) error {
	item.CreatedAt = time.Now()
	return d.Db.Create(item).Error
}

func (d *DailySummaryRepository) DeleteReport(date time.Time) error {
	dateFormmatedmin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dateFormmatedMax := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	return d.Db.Where("date BETWEEN ? and ?", dateFormmatedmin, dateFormmatedMax).Delete(&entities.DailySummary{}).Error
}