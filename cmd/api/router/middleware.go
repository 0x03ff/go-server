package router

import (
	"log"
	"net/http"
	"net/url"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func JWTMiddleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Extract user_id from URL param
			user_id := chi.URLParam(r, "user_id")
			if user_id == "" {
				log.Printf("User ID not found in URL path: %s", r.URL.Path)
				redirectWithError(w, r, "invalid_url", "/login")
				return
			}

			// Get the token from the cookie
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				log.Printf("Token cookie not found: %v", err)
				redirectWithError(w, r, "missing_token", "/login")
				return
			}

			tokenStr := tokenCookie.Value
			if tokenStr == "" {
				log.Println("Empty token in cookie")
				redirectWithError(w, r, "empty_token", "/login")
				return
			}

			// Verify the token
			systemRepo := repositories.NewKeysRepository(db)
			tokenObj, err := utils.VerifyToken(tokenStr, systemRepo)
			if err != nil {
				log.Printf("Token verification failed: %v", err)
				redirectWithError(w, r, "invalid_token", "/login")
				return
			}

			// Extract claims
			claims, ok := tokenObj.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Invalid token claims")
				redirectWithError(w, r, "invalid_claims", "/login")
				return
			}

			userIdClaim, ok := claims["user_id"].(string)
			if !ok {
				log.Println("user_id claim missing or invalid")
				redirectWithError(w, r, "missing_user_id", "/login")
				return
			}

			// Check if user_id in URL matches user_id in token
			if userIdClaim != user_id {
				log.Println("User ID mismatch between URL and token")
				redirectWithError(w, r, "id_mismatch", "/login")
				return
			}

			// All good, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to redirect with error message
func redirectWithError(w http.ResponseWriter, r *http.Request, errorCode, path string) {
	query := url.Values{}
	query.Set("error", errorCode)
	
	// Preserve the original path for redirect after login
	if r.URL.Path != "/login" {
		query.Set("redirect", r.URL.Path)
	}
	
	http.Redirect(w, r, path+"?"+query.Encode(), http.StatusSeeOther)
}
