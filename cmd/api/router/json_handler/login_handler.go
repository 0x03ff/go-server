package json_handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"log"
	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
)

func (h *JsonHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// ====== NEW: Rate limiting and brute-force protection ======
	clientID := h.getClientIdentifier(r)
	
	// Check if account is locked out
	if h.isLockedOut(clientID) {
		utils.SendError(w, http.StatusTooManyRequests, "Account locked due to multiple failed attempts. Please try again later.")
		return
	}
	
	// Apply progressive delay based on failed attempts
	h.applyProgressiveDelay(clientID)
	
	// Track the time of this login attempt
	h.updateLastLoginTime(clientID)
	// ====== END OF NEW RATE LIMITING SECTION ======

	csrf_err := utils.VerifyCSRFtoken(w, r)
	if csrf_err != nil {
		utils.SendError(w, http.StatusForbidden, "CSRF_TOKEN_INVALID")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// ====== NEW: Increment failed attempts on invalid request ======
		h.incrementFailedAttempts(clientID)
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = utils.ValidateUserInput("User ID", user.Username, 6, 20)
	if err != nil {
		// ====== NEW: Increment failed attempts on validation failure ======
		h.incrementFailedAttempts(clientID)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Password", user.Password, 8, 20)
	if err != nil {
		// ====== NEW: Increment failed attempts on validation failure ======
		h.incrementFailedAttempts(clientID)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Recover key", user.Recover, 6, 20)
	if err != nil {
		// ====== NEW: Increment failed attempts on validation failure ======
		h.incrementFailedAttempts(clientID)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	userRepo := repositories.NewUsersRepository(h.dbPool)
	loggedInUser, err := userRepo.Login(r.Context(), user.Username, user.Password, user.Recover)
	if err != nil {
		// ====== NEW: Increment failed attempts on login failure ======
		h.incrementFailedAttempts(clientID)
		
		// ====== NEW: Log suspicious activity ======
		h.logSuspiciousActivity(clientID, user.Username, "Invalid credentials")
		
		utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// ====== NEW: Clear failed attempts on successful login ======
	h.clearFailedAttempts(clientID)
	
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

// ====== NEW: Helper methods for rate limiting and brute-force protection ======

// getClientIdentifier determines the client identifier for rate limiting
// Can be based on IP address or username (or both)
func (h *JsonHandlers) getClientIdentifier(r *http.Request) string {
	// Get client IP (remove port if present)
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	
	// For more robust protection, you could also include the username if available
	// But be careful not to leak username information in logs
	return ip
}

// isLockedOut checks if the client is currently locked out
func (h *JsonHandlers) isLockedOut(clientID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if lockoutTime, exists := h.lockoutTimes[clientID]; exists {
		if time.Now().Before(lockoutTime) {
			return true
		}
		// Clean up expired lockout
		delete(h.lockoutTimes, clientID)
	}
	return false
}

// applyProgressiveDelay applies increasing delays based on failed attempts
func (h *JsonHandlers) applyProgressiveDelay(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	attempts := h.failedAttempts[clientID]
	if attempts > 0 {
		// Progressive delay: 1s, 2s, 4s, 8s, etc. (max 30 seconds)
		delay := time.Duration(1 << uint(attempts-1))
		if delay > 30 {
			delay = 30
		}
		
		// Log the delay being applied
		log.Printf("Applying %d second delay for client %s (attempt %d)", delay, clientID, attempts)
		
		time.Sleep(delay * time.Second)
	}
}

// incrementFailedAttempts increments the failed attempts counter and checks lockout
func (h *JsonHandlers) incrementFailedAttempts(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.failedAttempts[clientID]++
	attempts := h.failedAttempts[clientID]
	
	// Log the failed attempt
	log.Printf("Failed login attempt %d for client %s", attempts, clientID)
	
	// Apply lockout after 5 failed attempts
	if attempts >= 5 {
		lockoutDuration := time.Duration(15+5*(attempts-5)) * time.Minute
		h.lockoutTimes[clientID] = time.Now().Add(lockoutDuration)
		
		// Log the lockout
		log.Printf("Client %s locked out for %s due to %d failed attempts", 
			clientID, lockoutDuration, attempts)
	}
}

// clearFailedAttempts resets the failed attempts counter for a client
func (h *JsonHandlers) clearFailedAttempts(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	delete(h.failedAttempts, clientID)
	delete(h.lockoutTimes, clientID)
	
	// Log successful login
	log.Printf("Successful login for client %s", clientID)
}

// updateLastLoginTime records the time of the current login attempt
func (h *JsonHandlers) updateLastLoginTime(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.lastLoginTimes[clientID] = time.Now()
}

// logSuspiciousActivity logs suspicious login activity for monitoring
func (h *JsonHandlers) logSuspiciousActivity(clientID, username, reason string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	attempts := h.failedAttempts[clientID]
	
	// Log to system log
	log.Printf("Suspicious login activity: client=%s, username=%s, reason=%s, attempts=%d",
		clientID, username, reason, attempts)
	
	// You could also send alerts for high-risk activity
	if attempts >= 3 {
		// In a real system, you might send an email alert here
		log.Printf("High-risk login activity detected for username=%s", username)
	}
}