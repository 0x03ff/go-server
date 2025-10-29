package security

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Use Handler as the type name
type Handler struct {
	dbPool *pgxpool.Pool
}

// NewHandler accepts dbPool
func NewHandler(dbPool *pgxpool.Pool) *Handler {
	return &Handler{
		dbPool: dbPool,
	}
}

func (h *Handler) SecurityHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Security protocol overhead endpoint is active",
		"status":  "success",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
