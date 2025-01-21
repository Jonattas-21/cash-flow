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
	GetDailySummary(date time.Time) (entities.DailySummary, error)
	GenerateReport(date time.Time) (entities.DailySummary, error)
}

type DailySummaryUseCase struct {
	Repository         interfaces.DailySummaryRepository
	TransactionUseCase Transaction
	CashinCashoutUrl   string
	Rdb                *redis.Client //retirar o redis daqui
}

func (d *DailySummaryUseCase) GetDailySummary(date time.Time) (entities.DailySummary, error) {
	cacheKey := fmt.Sprintf("daily_summary_%s", date.Format("2006-01-02"))

	_, err := d.Rdb.Ping().Result()
	var summary entities.DailySummary

	// if the redis is up, let's try to get the data from cache
	if err == nil {
		log.Println("Trying to get data from cache with the key: ", cacheKey)
		val, err := d.Rdb.Get(cacheKey).Result()
		log.Println("Value from cache: ", val)
		if err == nil && val != "" {
			// Is not in the cache, lets find it in the DB
			err = json.Unmarshal([]byte(val), &summary)
			if err != nil {
				log.Println("Error to unmarshal the data from cache: ", err)
				return summary, err
			}

			log.Println("Found in cache for date: ", date)
			return summary, nil
		}
	} else {
		log.Println("Error to retrieve data from cache or reach the redis: ", err)
	}

	//It's today, let's delete the partial report
	today := time.Now()
	if date.Day() == today.Day() && date.Month() == today.Month() && date.Year() == today.Year() {
		err = d.Repository.DeleteReport(date)
		if err != nil {
			log.Println("Error to delete the report: ", err)
			return summary, err
		}
	}

	// if it's not in the cache or error in cache, let's find in the db
	log.Println("Not found in cache, let's find in the db with the date: ", date)
	summary, err = d.Repository.GetReport(date)
	if err != nil {
		return summary, err
	}

	if summary.Total == 0 {
		// if it not find in db, let's generate the report
		log.Println("Not found in db, let's generate the report for the date: ", date)
		summary, err = d.GenerateReport(date)
		if err != nil {
			return summary, err
		}
	}
	jsonSummary, _ := json.Marshal(summary)
	redisCommand := d.Rdb.Set(cacheKey, jsonSummary, time.Duration(2) * time.Minute)
	if redisCommand.Err() != nil {
		log.Println("Error to save the data in cache: ", redisCommand.Err())
	}

	log.Println("Summary: ", summary)
	return summary, nil
}

func (d *DailySummaryUseCase) GenerateReport(date time.Time) (entities.DailySummary, error) {
	transactions, err := d.TransactionUseCase.FindTransactions(date)
	var dailySymmary entities.DailySummary

	if err != nil {
		log.Println("Error to get transactions by date: ", err)
		return dailySymmary, err
	}

	
	dailySymmary.Date = date
	dailySymmary.CreatedAt = time.Now()
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
		return dailySymmary, err
	}

	return dailySymmary, nil
}
