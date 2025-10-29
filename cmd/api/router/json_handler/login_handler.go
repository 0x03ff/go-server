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
	loggedInUser, err := userRepo.Login(r.Context(), user.Username, user.Password, user.Recover)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	user_PrivateKey, err := userRepo.GetUserECDHPrivateKey(r.Context(), loggedInUser)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid client privary key")
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
		Secure:   true, // Set to true to using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 4, // 4 hour session
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)

	// Fetch the public key from the system repository
	var systemKey models.SystemKey

	publicKeyPem, err := systemRepo.GetECDHPublicKey(r.Context(), &systemKey)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get obtain public key")
		return

	}



	// Send a JSON response with user_id and token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":        loggedInUser.ID,
		"client_private": user_PrivateKey,
		"server_public":  publicKeyPem,
	})

}
