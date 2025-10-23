// router/json_handler/router.go
package json_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JsonHandlers struct{ dbPool *pgxpool.Pool }

func NewHandlers(dbPool *pgxpool.Pool) *JsonHandlers {
	return &JsonHandlers{dbPool: dbPool}
}

type App interface {
	GetJsonHandlers() *JsonHandlers
}

func SetupJsonRoutes(r chi.Router, app App) {
	handlers := app.GetJsonHandlers()

	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/register", handlers.RegisterHandler)

	r.Post("/api/upload_file/{user_id}", handlers.UploadFileHandler)


}
