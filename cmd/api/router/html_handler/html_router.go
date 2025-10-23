package html_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HtmlHandlers struct {
	dbPool *pgxpool.Pool
}

func NewHandlers(dbPool *pgxpool.Pool) *HtmlHandlers {
	return &HtmlHandlers{dbPool: dbPool}
}

type App interface {
	GetHtmlHandlers() *HtmlHandlers
}

func SetupHtmlRoutes(r chi.Router, app App) {
	handlers := app.GetHtmlHandlers()

	r.Get("/", handlers.IndexHandler)
	r.Get("/login", handlers.LoginHandler)
	r.Get("/register", handlers.RegisterHandler)


	r.Get("/transfer/{user_id}", handlers.TransferHandler)

	r.Get("/home/{user_id}", handlers.HomeHandler)
	r.Get("/file_download/{user_id}", handlers.FileDownloadHandler)
	r.Get("/file_upload/{user_id}", handlers.FileUploadHandler)
	r.Get("/logout/{user_id}", handlers.LogoutHandler)
}
