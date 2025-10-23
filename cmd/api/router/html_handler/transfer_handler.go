package html_handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-chi/chi/v5"
)

func (h *HtmlHandlers) TransferHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user_id from the URL path parameter
	user_id := chi.URLParam(r, "user_id")

	if user_id == "" {
		http.Error(w, "User ID not found in URL path", http.StatusBadRequest)
		return
	}

	// Get the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header not found.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract the token from the Authorization header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("Invalid authorization header format.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenStr := parts[1]

	if tokenStr == "" {
		log.Println("Invalid token format.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	systemRepo := repositories.NewKeysRepository(h.dbPool)
	// Verify the token
	tokenObj, err := utils.VerifyToken(tokenStr, systemRepo)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract claims from the token
	claims := tokenObj.Claims.(jwt.MapClaims)
	userIdClaim := claims["user_id"].(string)
	userNameClaim := claims["user_name"].(string)

	// Check if the user_id in the URL matches the user_id in the token
	if userIdClaim != user_id {
		log.Println("User ID mismatch.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load and parse the template
	tmpl, err := template.ParseFiles("web/html/transfer.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the user's name and user_id
	data := struct {
		UserName string
		UserID   string
	}{
		UserName: userNameClaim,
		UserID:   user_id,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
