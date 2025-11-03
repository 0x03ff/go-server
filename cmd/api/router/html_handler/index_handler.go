// router/html_handler/router.go
package html_handler

import (
	"net/http"
	"time"
)


func (h *HtmlHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
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
	
    http.ServeFile(w, r, "web/html/index.html")
}