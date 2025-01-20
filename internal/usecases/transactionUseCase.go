package usecases

import (
	"fmt"
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/domain/interfaces"
	"log"
	"time"

	"github.com/rs/xid"
)

type Transaction interface {
	SaveTransaction(transaction entities.Transaction) (entities.Transaction, []string, error)
	FindTransactions(date time.Time) ([]entities.Transaction, error)
	FindAllTransactions() ([]entities.Transaction, error)
	DeleteTransaction(id string) error
	UpdateTransaction(id string, item entities.Transaction) error
}

type TransactionUseCase struct {
	Repository interfaces.TransactionRepository
}

func (t *TransactionUseCase) SaveTransaction(transaction entities.Transaction) (entities.Transaction, []string, error) {
	errorsValidation := []string{}

	//validations
	if transaction.Amount <= 0 {
		log.Println(fmt.Println("Amount is required and must be greater than 0, received: ", transaction.Amount))	
		errorsValidation = append(errorsValidation, "Amount is required and must be greater than 0")
	}

	if transaction.Type == "" || (transaction.Type != "credit" && transaction.Type != "debit") {
		log.Println(fmt.Println("Type is required, received:", transaction.Type))
		errorsValidation = append(errorsValidation, "Transaction type is required and must be credit or debit")
	}

	if transaction.Date == (time.Time{}) || !transaction.Date.Before(time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1)) {
		log.Println(fmt.Println("Date is required and must be less or equal than today, received:", transaction.Date))
		errorsValidation = append(errorsValidation, "Date is required and must be less or equal than today")
	}

	if len(errorsValidation) > 0 {
		return transaction, errorsValidation, nil
	} else {

		transaction.CreatedAt = time.Now()
		transaction.ID = xid.New().String()

		//save database
		err := t.Repository.Save(transaction)
		if err != nil {
			log.Println("Error to save transaction on services: ", err)
			return transaction, nil, err
		}

		return transaction, errorsValidation, nil
	}
}

func (t *TransactionUseCase) FindAllTransactions() ([]entities.Transaction, error) {
	transactions, err := t.Repository.FindAll()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionUseCase) FindTransactions(date time.Time) ([]entities.Transaction, error) {
	transactions, err := t.Repository.FindByDay(date)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

//obsolete???
// func (d *TransactionUseCase) GetTransactionsByDate(baseURL string, date string) ([]Transaction, error) {
// 	u, err := url.Parse(baseURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	q := u.Query()
// 	q.Set("date", date)
// 	u.RawQuery = q.Encode()

// 	resp, err := http.Get(u.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
// 	}

// 	var transactions []Transaction
// 	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
// 		return nil, err
// 	}

// 	return transactions, nil
// }

func (t *TransactionUseCase) DeleteTransaction(id string) error {
	err := t.Repository.DeleteTransaction(id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionUseCase) UpdateTransaction(id string, item entities.Transaction) error {
	err := t.Repository.UpdateTransaction(id, item)
	if err != nil {
		return err
	}
	return nil
}
