package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/domain/interfaces"

	"github.com/go-redis/redis"
)

// todo: look over if it's necessary to have both methods to search
type DailySummary interface {
	GetDailySummary(date time.Time) (*entities.DailySummary, error)
	GenerateReport(date time.Time) (*entities.DailySummary, error)
}

type DailySummaryUseCase struct {
	Repository         interfaces.DailySummaryRepository
	TransactionUseCase Transaction
	CashinCashoutUrl   string
	Rdb                *redis.Client //retirar o redis daqui
}

func (d *DailySummaryUseCase) GetDailySummary(date time.Time) (*entities.DailySummary, error) {
	cacheKey := fmt.Sprintf("daily_summary:%s", date.Format("2006-01-02"))

	_, err := d.Rdb.Ping().Result()
	var summary *entities.DailySummary

	// if the redis is up, let's try to get the data from cache
	if err == nil {
		val, err := d.Rdb.Get(cacheKey).Result()
		if err == redis.Nil {
			// Is not in the cache, lets find it in the DB
			err = json.Unmarshal([]byte(val), &summary)
			if err != nil {
				return nil, err
			}

			return summary, nil
		}
	} else {
		log.Println("Error to retrieve data from cache or reach the redis: ", err)
	}

	// if it's not in the cache or error in cache, let's find in the db
	summary, err = d.Repository.GetReport(date)
	if err != nil {
		return nil, err
	}

	if summary == nil {
		// if it not find in db, let's generate the report
		summary, err = d.GenerateReport(date)
		if err != nil {
			return nil, err
		}
	}

	return summary, nil
}

func (d *DailySummaryUseCase) GenerateReport(date time.Time) (*entities.DailySummary, error) {
	transactions, err := d.TransactionUseCase.FindTransactions(date)

	if err != nil {
		log.Println("Error to get transactions by date: ", err)
		return nil, err
	}

	var dailySymmary entities.DailySummary
	dailySymmary.Date = date
	for _, t := range transactions {

		if t.Type == "credit" {
			dailySymmary.Credit += t.Amount
		} else {
			dailySymmary.Debit += t.Amount
		}

		dailySymmary.Total = dailySymmary.Credit - dailySymmary.Debit
	}

	// if the date is minor then today, i will calculate the report and set closed
	if (date.Year() < time.Now().Year()) ||
		(date.Year() == time.Now().Year() && date.Month() < time.Now().Month()) ||
		(date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() < time.Now().Day()) {

		dailySymmary.Status = "closed"
	} else {
		dailySymmary.Status = "partial"
	}

	// save in db
	err = d.Repository.SaveReport(dailySymmary)
	if err != nil {
		return nil, err
	}

	return &dailySymmary, nil
}
