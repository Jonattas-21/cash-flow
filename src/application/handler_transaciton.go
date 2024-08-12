package application

import (
	"cash-flow/src/application/dto"
	"cash-flow/src/domain/transaction"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HandlerTransaction struct {
	TransactionUseCase transaction.ITransactionUseCase
}

func (h *HandlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction transaction.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Fatalln("Error to decode transaction: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.CreatedBy = r.Context().Value("email").(string)

	result, validations, err := h.TransactionUseCase.SaveTransaction(transaction)
	if err != nil {
		log.Fatalln("Error to save transaction: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.Response{
		Timestamp: time.Now(),
		Object:    result,
		Messages:  validations,
	}

	log.Println("Transaction saved: ", result)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *HandlerTransaction) GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		http.Error(w, "Missing 'date' query parameter", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	transactions, err := h.TransactionUseCase.FindTransactions(date)
	if err != nil {
		log.Fatalln("Error to get transactions: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
