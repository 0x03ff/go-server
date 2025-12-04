package web_server

import (
	"net/http"

	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server/server"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool" // Add this import
)

// WebServerHandlers holds all the handlers for the resilience API
type WebServerHandlers struct {
	Server *server.Handler
}

// NewWebServerHandlers creates and returns all resilience handlers WITH dbPool
func NewWebServerHandlers(dbPool *pgxpool.Pool) *WebServerHandlers {
	return &WebServerHandlers{
		Server: server.NewHandler(dbPool),
	}
}

// App interface for accessing resilience handlers
type App interface {
	GetWebServerHandlers() *WebServerHandlers
}

// SetupResilienceRoutes configures all resilience API routes
func SetupResilienceRoutes(r chi.Router, app App) {
	handlers := app.GetWebServerHandlers()

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/image/storage_icon.ico", http.StatusMovedPermanently)
	})
	// Web Server Hardening Baseline
	r.Route("/api/web_server", func(r chi.Router) {
		r.Get("/", handlers.Server.WebServerHandler)
	})

}
