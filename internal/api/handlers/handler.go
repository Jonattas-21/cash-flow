package handlers

import (
	"net/http"
)

// @Summary Check if the service is running
// @Description Check if the service is running
// @Produce json
// @Success 200 {string} string "Health Check!"
// @Router / [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	parametroString := r.URL.Query().Get("CampaignName")
	w.Write([]byte("Health Check! " + parametroString))
}
