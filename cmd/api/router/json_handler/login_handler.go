package json_handler

import (
	"encoding/json"
	"net/http"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
)

func (h *JsonHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        utils.SendError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    userRepo := repositories.NewUsersRepository(h.dbPool)
    loggedInUser, err := userRepo.Login(r.Context(), user.Username, user.Password)
    if err != nil {
        utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    systemRepo := repositories.NewKeysRepository(h.dbPool)

    token, err := utils.GenerateToken(r.Context(), loggedInUser.ID, loggedInUser.Username, systemRepo)
    if err != nil {
        utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
        return
    }

    // Create a new cookie
    cookie := &http.Cookie{
        Name:     "token",
        Value:    token,
        HttpOnly: true,
        Path:     "/",
        Secure:   r.TLS != nil, // Set to true if using HTTPS
        SameSite: http.SameSiteLaxMode,
    }

    // Set the cookie in the response
    http.SetCookie(w, cookie)

    // Send a JSON response with user_id and token
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "user_id": loggedInUser.ID,
        "token":   token,
    })
}
