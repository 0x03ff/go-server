package json_handler

import (
	"encoding/json"
	"net/http"

	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (h *JsonHandlers) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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

 

    // Return a success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}
