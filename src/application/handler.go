package application

import (
	"cash-flow/src/domain/transaction"
	"net/http"
)

type Handler struct {
	TransactionUseCase transaction.IUseCaseTransaction
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	parametroString := r.URL.Query().Get("CampaignName")
	w.Write([]byte("Helth Check! " + parametroString))
}
