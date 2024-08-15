package transaction

import (
	"log"
	"time"

	"github.com/rs/xid"
)

type ITransactionUseCase interface {
	SaveTransaction(transaction Transaction) (Transaction, []string, error)
	FindTransactions(date time.Time) ([]Transaction, error)
}

type TransactionUseCase struct {
	Repository ITransactionRepository
}

func (t *TransactionUseCase) SaveTransaction(transaction Transaction) (Transaction, []string, error) {
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

func (t *TransactionUseCase) FindTransactions(date time.Time) ([]Transaction, error) {
	transactions, err := t.Repository.FindByDay(date)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
