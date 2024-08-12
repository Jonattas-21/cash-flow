package transaction

import (
	"cash-flow/src/infrastructure/database"
	"log"
	"time"
)

type IUseCaseTransaction interface {
	SaveTransaction(transaction Transaction) (Transaction, []string, error)
	FindTransactions(date time.Time) ([]Transaction, error)
}

type TransactionUseCase struct {
	IRepository database.IRepository
}

func (t *TransactionUseCase) SaveTransaction(transaction Transaction) (Transaction, []string, error) {
	transaction.CreatedAt = time.Now()
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

	//save database
	err := t.IRepository.Save(transaction)
	if err != nil {
		log.Fatalln("Error to save transaction on services: ", err)
		return transaction, nil, err
	}

	return transaction, errorsValidation, nil
}

func (t *TransactionUseCase) FindTransactions(date time.Time) ([]Transaction, error) {
	transactions, err := t.IRepository.FindByDay(date)
	if err != nil {
		return nil, err
	}

	var result []Transaction
	for _, item := range transactions {
		if v, ok := item.(Transaction); ok {
			result = append(result, v)
		} else {
			log.Printf("Error to convert transaction to Transaction %v", item)
		}
	}

	return result, nil
}
