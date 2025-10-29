package router

import (
	"log"
	"net/http"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func JWTMiddleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Incoming URL: %s", r.URL.Path)


			
			// Extract user_id from URL param
			user_id := chi.URLParam(r, "user_id")
			if user_id == "" {
				log.Printf("User ID not found in URL path: %s", r.URL.Path)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Get the token from the cookie
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				log.Printf("Token cookie not found: %v", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			tokenStr := tokenCookie.Value
			if tokenStr == "" {
				log.Println("Empty token in cookie")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Verify the token
			systemRepo := repositories.NewKeysRepository(db)
			tokenObj, err := utils.VerifyToken(tokenStr, systemRepo)
			if err != nil {
				log.Printf("Token verification failed: %v", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Extract claims
			claims, ok := tokenObj.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Invalid token claims")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userIdClaim, ok := claims["user_id"].(string)
			if !ok {
				log.Println("user_id claim missing or invalid")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}


			// Check if user_id in URL matches user_id in token
			if userIdClaim != user_id {
				log.Println("User ID mismatch between URL and token")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// All good, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}
