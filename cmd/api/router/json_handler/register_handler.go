package json_handler

import (
	"encoding/json"
	"net/http"


	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
)

func (h *JsonHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	clientID := h.GetClientIdentifier(r)
	
	// Check if client is locked out (this now includes registration locks)
	if h.IsLockedOut(clientID) {
		utils.SendError(w, http.StatusTooManyRequests, "Account creation locked due to multiple accounts from same IP. Please try again later.")
		return
	}
	
	// Apply progressive delay based on failed attempts
	h.ApplyProgressiveDelay(clientID, "register")
	
	// Track the time of this registration attempt
	h.UpdateLastLoginTime(clientID)

	csrf_err := utils.VerifyCSRFtoken(w, r)
	if csrf_err != nil {
		// Increment failed attempts on CSRF failure
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusForbidden, "CSRF_TOKEN_INVALID")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Increment failed attempts on invalid request
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate user data
	if user.Username == "" || user.Password == "" {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	err = utils.ValidateUserInput("User ID", user.Username, 6, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Password", user.Password, 8, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}	
	err = utils.ValidateUserInput("Recover key", user.Recover, 6, 20)
	if err != nil {
		// Increment failed attempts on validation failure
		h.IncrementFailedAttempts(clientID, "register")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	

	userRepo := repositories.NewUsersRepository(h.dbPool)
	err = userRepo.Create(r.Context(), &user)

	if err != nil {
		// Increment failed attempts on registration failure
		h.IncrementFailedAttempts(clientID, "register")
		
		// Log suspicious activity
		h.LogSuspiciousActivity(clientID, user.Username, "Registration failed", "register")
		
		utils.SendError(w, http.StatusInternalServerError, "Registration failed. Please try again.")
		return
	}

	// Clear failed attempts on successful registration
	h.ClearFailedAttempts(clientID, user.Username, "register")
	
	// ====== NEW: Track successful registrations and apply IP lockout ======
	h.IncrementSuccessfulRegistrations(clientID,user.Username)
	
	// Check if this IP has created too many accounts
	if h.IsRegistrationThresholdReached(clientID) {
		h.LockIPForRegistration(clientID)
		h.LogSuspiciousActivity(clientID, user.Username, "Multiple accounts created from same IP", "register")
		utils.SendError(w, http.StatusTooManyRequests, "Account creation locked due to multiple accounts from same IP. Please try again later.")
		return
	}
	// ====== END OF NEW IP TRACKING SECTION ======

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
	})
}
