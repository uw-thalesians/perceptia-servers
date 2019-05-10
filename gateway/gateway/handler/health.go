package handler

import (
	"encoding/json"
	"net/http"
)

func (cx *Context) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type healthObj struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}

	healthStatus := healthObj{
		Name:    "Perceptia API Health Report",
		Version: cx.gatewayVersion,
		Status:  "ready",
	}
	w.WriteHeader(http.StatusOK)
	errWJ := json.NewEncoder(w).Encode(healthStatus)
	if errWJ != nil {
		cx.logError(r, errWJ, "trying to encode struct as json and write to response",
			"", http.StatusInternalServerError)
	}
	return
}
