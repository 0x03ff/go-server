package json_handler

import (
	"github.com/0x03ff/golang/cmd/api/router"
	"github.com/0x03ff/golang/cmd/api/router/json_handler/web_server"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// JsonHandlers holds all JSON API handlers
type JsonHandlers struct {
	dbPool    *pgxpool.Pool
	WebServer *web_server.WebServerHandlers
}

// NewHandlers creates and returns all JSON handlers
func NewHandlers(dbPool *pgxpool.Pool) *JsonHandlers {
	return &JsonHandlers{
		dbPool:    dbPool,
		WebServer: web_server.NewWebServerHandlers(dbPool), // Pass dbPool here
	}
}

// App interface for accessing JSON handlers
type App interface {
	GetJsonHandlers() *JsonHandlers
}

// GetWebServerHandlers implements the web_server.App interface
func (h *JsonHandlers) GetWebServerHandlers() *web_server.WebServerHandlers {
	return h.WebServer
}

// SetupJsonRoutes configures all JSON API routes
func SetupJsonRoutes(r chi.Router, app App) {
	handlers := app.GetJsonHandlers()

	// Setup resilience routes
	web_server.SetupResilienceRoutes(r, handlers)

	// Existing routes
	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/register", handlers.RegisterHandler)

	authMiddleware := router.JWTMiddleware(handlers.dbPool)

	r.With(authMiddleware).Get("/api/files/{user_id}", handlers.ObtainFileHandler)
	
	// Keep these as they are
	r.With(authMiddleware).Post("/api/upload_file/{user_id}", handlers.UploadFileHandler)
	r.With(authMiddleware).Get("/api/download_file/{user_id}/{file_id}", handlers.DownloadFileHandler)
}
