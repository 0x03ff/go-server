package html_handler

import (
	"html/template"
	"net/http"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (h *HtmlHandlers) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")

	if user_id == "" {
		http.Error(w, "User ID not found in URL path", http.StatusBadRequest)
		return
	}

	// Get the token from the cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token not found. Please log in again.", http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	if token == "" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	systemRepo := repositories.NewKeysRepository(h.dbPool)
	// Verify the token
	tokenObj, err := utils.VerifyToken(token, systemRepo)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract claims from the token
	claims := tokenObj.Claims.(jwt.MapClaims)
	userIdClaim, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Check if the user_id in the URL matches the user_id in the token
	if userIdClaim != user_id {
		http.Error(w, "User ID mismatch", http.StatusUnauthorized)
		return
	}

	// Render the template with just the user_id
	tmpl, err := template.ParseFiles("web/html/file_download.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"UserID": user_id,
	})

	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
