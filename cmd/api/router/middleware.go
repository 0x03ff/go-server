package router

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func JWTMiddleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ====== NEW: Determine if this is an API request ======
			isAPIRequest := strings.HasPrefix(r.URL.Path, "/api/")
			// ====== END OF NEW SECTION ======

			// Extract user_id from URL param
			user_id := chi.URLParam(r, "user_id")
			if user_id == "" {
				log.Printf("User ID not found in URL path: %s", r.URL.Path)
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusBadRequest, "invalid_url", "User ID not found in URL path")
				} else {
					redirectWithError(w, r, "invalid_url", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			// Get the token from the cookie
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				log.Printf("Token cookie not found: %v", err)
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "missing_token", "Session expired. Please log in again.")
				} else {
					redirectWithError(w, r, "missing_token", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			tokenStr := tokenCookie.Value
			if tokenStr == "" {
				log.Println("Empty token in cookie")
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "empty_token", "Invalid session. Please log in again.")
				} else {
					redirectWithError(w, r, "empty_token", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			// Verify the token
			systemRepo := repositories.NewKeysRepository(db)
			tokenObj, err := utils.VerifyToken(tokenStr, systemRepo)
			if err != nil {
				log.Printf("Token verification failed: %v", err)
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "invalid_token", "Invalid session token. Please log in again.")
				} else {
					redirectWithError(w, r, "invalid_token", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			// Extract claims
			claims, ok := tokenObj.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Invalid token claims")
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "invalid_claims", "Session data corrupted. Please log in again.")
				} else {
					redirectWithError(w, r, "invalid_claims", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			userIdClaim, ok := claims["user_id"].(string)
			if !ok {
				log.Println("user_id claim missing or invalid")
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "missing_user_id", "User ID missing from session. Please log in again.")
				} else {
					redirectWithError(w, r, "missing_user_id", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
				return
			}

			// Check if user_id in URL matches user_id in token
			if userIdClaim != user_id {
				log.Println("User ID mismatch between URL and token")
				// ====== MODIFIED: Handle API vs Web requests ======
				if isAPIRequest {
					utils.SendJSONError(w, http.StatusUnauthorized, "id_mismatch", "User ID mismatch. Please log in again.")
				} else {
					redirectWithError(w, r, "id_mismatch", "/login")
				}
				// ====== END OF MODIFIED SECTION ======
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
