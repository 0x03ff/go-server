package html_handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-chi/chi/v5"
)

func (h *HtmlHandlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
    user_id := chi.URLParam(r, "user_id")

    if user_id == "" {
        http.Error(w, "User ID not found in URL path", http.StatusBadRequest)
        return
    }

    log.Printf("User ID from URL: %s", user_id)

    // Get the token from the cookie
    tokenCookie, err := r.Cookie("token")
    if err != nil {
        log.Printf("Token cookie not found: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    tokenStr := tokenCookie.Value

    if tokenStr == "" {
        log.Println("Invalid token format.")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    log.Printf("Token from cookie: %s", tokenStr)

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

    log.Printf("User ID from token: %s", userIdClaim)
    log.Printf("User Name from token: %s", userNameClaim)

    // Check if the user_id in the URL matches the user_id in the token
    if userIdClaim != user_id {
        log.Println("User ID mismatch.")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Load and parse the template
    tmpl, err := template.ParseFiles("web/html/home.html")
    if err != nil {
        http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Render the template with the user's name and user_id
    data := struct {
        UserName string
        UserID   string
    }{
        UserName: userNameClaim,
        UserID:   user_id,
    }

    log.Printf("Rendering template with data: %+v", data)

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
        return
    }
}
