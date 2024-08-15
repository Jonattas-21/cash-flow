package handlers

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	parametroString := r.URL.Query().Get("CampaignName")
	w.Write([]byte("Health Check! " + parametroString))
}
