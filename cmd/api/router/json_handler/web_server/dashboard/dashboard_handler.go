package dashboard

import (
	"net/http"
	"encoding/json"
)

// Use Handler as the type name
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Dashboard endpoint is active",
		"status":  "success",
	}
	
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
