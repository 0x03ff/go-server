// router/html_handler/router.go
package html_handler

import (
	"net/http"
	"text/template"
	"time"

	"github.com/0x03ff/golang/utils"
)



func (h *HtmlHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	// Generate CSRF token
	csrfToken := utils.GenerateCSRFToken()
	
	// Set CSRF token in cookie (SameSite=Strict for login protection)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Secure:   r.TLS != nil,  // Must be true in production
		HttpOnly: false, // Needed for JavaScript access
		SameSite: http.SameSiteStrictMode,
		MaxAge:   300,   // 5 minutes validity
	})

	// Serve template with CSRF token
	tmpl := template.Must(template.ParseFiles("web/html/login.html"))
	tmpl.Execute(w, map[string]string{"CSRFToken": csrfToken})
}