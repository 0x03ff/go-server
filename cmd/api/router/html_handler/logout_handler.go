package html_handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *HtmlHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")

	if user_id == "" {
		http.Error(w, "User ID not found in URL path", http.StatusBadRequest)
		return
	}


	// Clear the JWT token from the cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: r.TLS != nil,
		Secure:   r.TLS != nil, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // Set the expiration date to a time in the past
	}
	http.SetCookie(w, cookie)


	// Load and parse the logout template
	tmpl, err := template.ParseFiles("web/html/logout.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the logout template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
