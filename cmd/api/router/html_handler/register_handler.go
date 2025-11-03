package html_handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/0x03ff/golang/utils"
)

func (h *HtmlHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// Clear the JWT token from the cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // Set the expiration date to a time in the past
	}
	http.SetCookie(w, cookie)
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
		MaxAge:   300, // 5 minutes validity
	})

	// Serve template with CSRF token
	tmpl := template.Must(template.ParseFiles("web/html/register.html"))
	tmpl.Execute(w, map[string]string{"CSRFToken": csrfToken})
}
