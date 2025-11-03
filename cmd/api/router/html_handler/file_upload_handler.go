package html_handler

import (
	"html/template"
	"net/http"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-chi/chi/v5"
)

func (h *HtmlHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
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
	userIdClaim := claims["user_id"].(string)

	// Check if the user_id in the URL matches the user_id in the token
	if userIdClaim != user_id {
		http.Error(w, "User ID mismatch", http.StatusUnauthorized)
		return
	}


	// Generate CSRF token
	csrfToken := utils.GenerateCSRFToken()
	
	// Set CSRF token in cookie (SameSite=Strict for login protection)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Secure:   true,  // Must be true in production
		HttpOnly: false, // Needed for JavaScript access
		SameSite: http.SameSiteStrictMode,
		MaxAge:   300,   // 5 minutes validity
	})	

	// Example: Render a template with the user_id and public key
	tmpl, err := template.ParseFiles("web/html/file_upload.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"UserID":     user_id,
		"CSRFToken": csrfToken,
	})


	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
