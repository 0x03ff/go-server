package web_server

import (
	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server/firewall"
	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server/security"
	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server/server"
	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server/tcp"
	"github.com/jackc/pgx/v5/pgxpool" // Add this import
	"github.com/go-chi/chi/v5"
)

// WebServerHandlers holds all the handlers for the resilience API
type WebServerHandlers struct {
	Server    *server.Handler
	Firewall  *firewall.Handler
	TCP       *tcp.Handler
	Security  *security.Handler
}

// NewWebServerHandlers creates and returns all resilience handlers WITH dbPool
func NewWebServerHandlers(dbPool *pgxpool.Pool) *WebServerHandlers {
	return &WebServerHandlers{
		Server:    server.NewHandler(dbPool),
		Firewall:  firewall.NewHandler(dbPool),
		TCP:       tcp.NewHandler(dbPool),
		Security:  security.NewHandler(dbPool),
	}
}

// App interface for accessing resilience handlers
type App interface {
	GetWebServerHandlers() *WebServerHandlers
}

// SetupResilienceRoutes configures all resilience API routes
func SetupResilienceRoutes(r chi.Router, app App) {
	handlers := app.GetWebServerHandlers()
	
	// Web Server Hardening Baseline (Researcher A)
	r.Route("/api/web_server", func(r chi.Router) {
		r.Get("/", handlers.Server.WebServerHandler)
		r.Post("/tests", handlers.Server.WebServerHandler)
	})
	
	// Firewall Rule Efficacy (Researcher B)
	r.Route("/api/firewall", func(r chi.Router) {
		r.Get("/", handlers.Firewall.FirewallHandler)
		r.Post("/tests", handlers.Firewall.FirewallHandler)
	})
	
	// TCP Stack Optimization (Researcher C)
	r.Route("/api/tcp", func(r chi.Router) {
		r.Get("/", handlers.TCP.TCPHandler)
		r.Post("/tests", handlers.TCP.TCPHandler)
	})
	
	// Security Protocol Overhead (Researcher D)
	r.Route("/api/security", func(r chi.Router) {
		r.Get("/", handlers.Security.SecurityHandler)
		r.Post("/tests", handlers.Security.SecurityHandler)
	})
}
