package dailySummary

import (
	"cash-flow/src/domain/transaction"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis"
)

type IDailySummaryUseCase interface {
	GetDailySummary(date time.Time) (*DailySummary, error)
	getTransactionsByDate(baseURL string, date string) ([]transaction.Transaction, error)
	getGeneratedReport(date time.Time) (*DailySummary, error)
	generateReport(date time.Time) (DailySummary, error)
}

type DailySummaryUseCase struct {
	Repository       IDailySummaryRepository
	CashinCashoutUrl string
	Rdb              *redis.Client
}

func (d *DailySummaryUseCase) GetDailySummary(date time.Time) (*DailySummary, error) {

	cacheKey := fmt.Sprintf("daily_summary:%s", date.Format("2006-01-02"))
	val, err := d.Rdb.Get(cacheKey).Result()
	var summary *DailySummary

	if err == redis.Nil {
		summary, err = d.getGeneratedReport(date)
		if err != nil {
			return nil, err
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

func (d *DailySummaryUseCase) getTransactionsByDate(baseURL string, date string) ([]transaction.Transaction, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("date", date)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	var transactions []transaction.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (d *DailySummaryUseCase) getGeneratedReport(date time.Time) (*DailySummary, error) {

	// if the date is minor then today, i will calculate the partial report
	var dateMinor bool
	if date.Year() < time.Now().Year() {
		dateMinor = true
	} else if date.Year() == time.Now().Year() && date.Month() < time.Now().Month() {
		dateMinor = true
	} else if date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() < time.Now().Day() {
		dateMinor = true
	}

	// if not, i will search in DB the closed report
	if dateMinor {
		result, err := d.Repository.GetReport(date)
		if err != nil {
			log.Fatal("Error to get report for date ", err, date)
			return nil, err
		}

		if result == nil {
			log.Println("Report not found for date, lets generate ", date)

			result, err := d.getGeneratedReport(date)
			if err != nil {
				log.Fatal("Error to generate report for date ", err, date)
				return nil, err
			}
			result.Status = "closed"
			err = d.Repository.SaveReport(*result)
			if err != nil {
				log.Fatal("Error to save report for date ", err, date)
				return nil, err
			}
		}
		return result, nil
	} else {
		result, err := d.getGeneratedReport(date)
		if err != nil {
			log.Fatal("Error to generate partial report for date ", err, date)
			return nil, err
		}

		result.Status = "partial"
		return result, nil
	}
}

func (d *DailySummaryUseCase) generateReport(date time.Time) (DailySummary, error) {
	transactions, err := d.getTransactionsByDate(d.CashinCashoutUrl, date.Format("2006-01-02"))

	if err != nil {
		log.Println("Error to get transactions by date: ", err)
		return DailySummary{}, err
	}

	var dailySymmary DailySummary
	dailySymmary.Date = date
	for _, t := range transactions {

		if t.Type == "credit" {
			dailySymmary.Credit += t.Amount
		} else {
			dailySymmary.Debit += t.Amount
		}

		dailySymmary.Total = dailySymmary.Credit - dailySymmary.Debit
	}

	return dailySymmary, nil
}
