package handlers

import (
	"encoding/json"
	"github.com/Jonattas-21/cash-flow/internal/api/dto"
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/usecases"
	"log"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
)

type HandlerTransaction struct {
	TransactionUseCase usecases.Transaction
}

// @Summary Create a new transaction
// @Description Create a new transaction, debit or credit
// @Produce json
// @Success 200 {transaction} single object
// @Router / [post]
func (h *HandlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction entities.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println("Error to decode transaction: ", err)
		http.Error(w, "Error to decode transaction: " + err.Error(), http.StatusBadRequest)
		return
	}
	if r.Context().Value("email") != nil {
		transaction.CreatedBy = r.Context().Value("email").(string)
	}

	result, validations, err := h.TransactionUseCase.SaveTransaction(transaction)
	if err != nil {
		log.Println("Error to save transaction: ", err)
		http.Error(w, "Error to save transaction: " + err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.Response{
		Timestamp: time.Now(),
		Messages:  validations,
	}

	if len(validations) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}


	response.Object = result
	log.Println("Transaction saved: ", result)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get transactions
// @Description Get all transactions by date
// @Produce json
// @Success 200 {transaction} list object
// @Router / [get]
func (h *HandlerTransaction) GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		transactions, err := h.TransactionUseCase.FindAllTransactions()
		if err != nil {
			log.Println("Error to get all transactions: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(transactions)
		
	} else{
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	
		transactions, err := h.TransactionUseCase.FindTransactions(date)
		if err != nil {
			log.Println("Error to get transactions: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		json.NewEncoder(w).Encode(transactions)
	}
}

// @Summary Update transaction
// @Description Update a transaction by id
// @Produce json
// @Success 200 {string} string "OK"
// @Router / [patch]
func (h *HandlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction entities.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println("Error to decode transaction: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if transaction.ID == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	err := h.TransactionUseCase.UpdateTransaction(transaction.ID, transaction)
	if err != nil {
		log.Println("Error to update transaction: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete transaction
// @Description Delete a transaction by id
// @Produce json
// @Success 200 {string} string "OK"
// @Router / [delete]
func (h *HandlerTransaction) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Println("ID: ", id)

	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	err := h.TransactionUseCase.DeleteTransaction(id)
	if err != nil {
		log.Println("Error to delete transaction: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}