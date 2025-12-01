package json_handler

import (
	"encoding/json"
	"net/http"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
)

func (h *JsonHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {

	clientID := h.GetClientIdentifier(r)
	
	// Check if account is locked out
	if h.IsLockedOut(clientID) {
		utils.SendError(w, http.StatusTooManyRequests, "Account locked due to multiple failed attempts. Please try again later.")
		return
	}
	
	// Apply progressive delay based on failed attempts
	h.ApplyProgressiveDelay(clientID,"login")
	
	// Track the time of this login attempt
	h.UpdateLastLoginTime(clientID)


	csrf_err := utils.VerifyCSRFtoken(w, r)
	if csrf_err != nil {
		utils.SendError(w, http.StatusForbidden, "CSRF_TOKEN_INVALID")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Increment failed attempts on invalid request
		h.IncrementFailedAttempts(clientID,"login")
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = utils.ValidateUserInput("User ID", user.Username, 6, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID,"login")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Password", user.Password, 8, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID,"login")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Recover key", user.Recover, 6, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID,"login")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	userRepo := repositories.NewUsersRepository(h.dbPool)
	loggedInUser, err := userRepo.Login(r.Context(), user.Username, user.Password, user.Recover)
	if err != nil {
		// Increment failed attempts on login failure
		h.IncrementFailedAttempts(clientID,"login")
		
		// Log suspicious activity
		h.LogSuspiciousActivity(clientID, user.Username, "Invalid credentials", "login")
		
		utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Clear failed attempts on successful login
	h.ClearFailedAttempts(clientID, "login")
	
	user_PrivateKey, err := userRepo.GetUserECDHPrivateKey(r.Context(), loggedInUser)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid client private key")
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
		Secure:   r.TLS != nil, // Set to true to using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 4, // 4 hour session
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)

	// Fetch the public key from the system repository
	var systemKey models.SystemKey
	publicKeyPem, err := systemRepo.GetECDHPublicKey(r.Context(), &systemKey)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get public key")
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
