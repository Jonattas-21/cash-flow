package application

import (
	"cash-flow/src/application/dto"
	"cash-flow/src/domain/transaction"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
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

// GetTransactions retorna todas as transações
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
