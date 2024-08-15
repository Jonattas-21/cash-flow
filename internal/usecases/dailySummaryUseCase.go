package usecases

import (
	"cash-flow/internal/domain/entities"
	"cash-flow/internal/domain/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis"
)

// todo rever se precisa de 2
type DailySummary interface {
	GetDailySummary(date time.Time) (*entities.DailySummary, error)
	GenerateReport(date time.Time) (*entities.DailySummary, error)
	getTransactionsByDate(date string) ([]entities.Transaction, error) //tem q ser DTO TODO
}

type DailySummaryUseCase struct {
	Repository       interfaces.DailySummaryRepository
	CashinCashoutUrl string
	Rdb              *redis.Client //retirar o redis daqui
}

func (d *DailySummaryUseCase) GetDailySummary(date time.Time) (*entities.DailySummary, error) {

	cacheKey := fmt.Sprintf("daily_summary:%s", date.Format("2006-01-02"))
	val, err := d.Rdb.Get(cacheKey).Result()
	var summary *entities.DailySummary

	if err == redis.Nil {
		// Is not int hte cache, lets find in the DB
		summary, err = d.Repository.GetReport(date)
		if err != nil {
			return nil, err
		}

		if summary == nil {
			// if not found in db, lets generate the report
			summary, err = d.GenerateReport(date)
			if err != nil {
				return nil, err
			}
		}
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (d *DailySummaryUseCase) GenerateReport(date time.Time) (*entities.DailySummary, error) {
	transactions, err := d.getTransactionsByDate(date.Format("2006-01-02"))

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

func (d *DailySummaryUseCase) getTransactionsByDate(date string) ([]entities.Transaction, error) {
	u, err := url.Parse(d.CashinCashoutUrl)
	if err != nil {
		return nil, err
	}

	u.Path = "/transactions"
	q := u.Query()
	q.Set("date", date)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var transactions []entities.Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
