package usecases

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/domain/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/xid"
)

type Transaction interface {
	SaveTransaction(transaction entities.Transaction) (entities.Transaction, []string, error)
	FindTransactions(date time.Time) ([]entities.Transaction, error)
}

type TransactionUseCase struct {
	Repository interfaces.TransactionRepository
}

func (t *TransactionUseCase) SaveTransaction(transaction entities.Transaction) (entities.Transaction, []string, error) {
	errorsValidation := []string{}

	//validations
	if transaction.Amount == 0 {
		log.Println("Amount is required")
		errorsValidation = append(errorsValidation, "Amount is required")
	}

	if transaction.Type == "" {
		log.Println("Type is required")
		errorsValidation = append(errorsValidation, "Transaction type is required and must be credit or debit")
	}

	transaction.CreatedAt = time.Now().Truncate(24 * time.Hour)
	transaction.ID = xid.New().String()

	//save database
	err := t.Repository.Save(transaction)
	if err != nil {
		log.Println("Error to save transaction on services: ", err)
		return transaction, nil, err
	}

	return transaction, errorsValidation, nil
}

func (t *TransactionUseCase) FindTransactions(date time.Time) ([]entities.Transaction, error) {
	transactions, err := t.Repository.FindByDay(date)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (d *TransactionUseCase) GetTransactionsByDate(baseURL string, date string) ([]Transaction, error) {
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

	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
