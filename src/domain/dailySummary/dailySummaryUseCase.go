package dailySummary

import (
	"cash-flow/src/domain/transaction"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type IDailySummaryUseCase interface {
	GenerateReport(date time.Time) (DailySummary, error)
	getTransactionsByDate(baseURL string, date string) ([]transaction.Transaction, error)
}

type DailySummaryUseCase struct {
	Repository       IDailySummaryRepository
	CashinCashoutUrl string
}

func (t *DailySummaryUseCase) GenerateReport(date time.Time) (DailySummary, error) {
	summaryMap := make(map[string]*DailySummary)
	transactions, err := t.getTransactionsByDate(t.CashinCashoutUrl, date.Format("2006-01-02"))

	if err != nil {
		log.Println("Error to get transactions by date: ", err)
		return DailySummary{}, err
	}

	for _, t := range transactions {
		date := t.CreatedAt.Format("2006-01-02")
		if summaryMap[date] == nil {
			summaryMap[date] = &DailySummary{
				Date: t.CreatedAt,
			}
		}
		if t.Type == "credit" {
			summaryMap[date].Credit += t.Amount
		} else {
			summaryMap[date].Debit += t.Amount
		}
		summaryMap[date].Total = summaryMap[date].Credit - summaryMap[date].Debit

		//todo implementar regra de negócio para status
	}

	var summaries []DailySummary
	for _, summary := range summaryMap {
		summaries = append(summaries, *summary)
	}

	return summaries[0], nil
}

func (t *DailySummaryUseCase) getTransactionsByDate(baseURL string, date string) ([]transaction.Transaction, error) {
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
