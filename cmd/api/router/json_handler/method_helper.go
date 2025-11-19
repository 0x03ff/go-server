package json_handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"


)

func generateRandomIP() string {
    return fmt.Sprintf("%d.%d.%d.%d", 
        rand.Intn(256), 
        rand.Intn(256), 
        rand.Intn(256), 
        rand.Intn(256))
}



// methods for rate limiting and brute-force protection

// GetClientIdentifier determines the client identifier for rate limiting
// Can be based on IP address
func (h *JsonHandlers) GetClientIdentifier(r *http.Request) string {
	// Get client IP (remove port if present)
	ip := r.RemoteAddr

	
	
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}


	if h.random_address{

	 ip = generateRandomIP()
	} 
	

	return ip
}

// IsLockedOut checks if the client is currently locked out
func (h *JsonHandlers) IsLockedOut(clientID string) bool {
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

// ApplyProgressiveDelay applies increasing delays based on failed attempts
func (h *JsonHandlers) ApplyProgressiveDelay(clientID string, action string) {
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
		log.Printf("Applying %d second delay for client " + action + " %s (attempt %d)", delay, clientID, attempts)
		
		time.Sleep(delay * time.Second)
	}
}

// IncrementFailedAttempts increments the failed attempts counter and checks lockout
func (h *JsonHandlers) IncrementFailedAttempts(clientID string, action string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.failedAttempts[clientID]++
	attempts := h.failedAttempts[clientID]
	
	// Log the failed attempt
	log.Printf("Failed " + action + " attempt %d for client %s", attempts, clientID)
	
	// Apply lockout after 5 failed attempts
	if attempts >= 5 {
		lockoutDuration := time.Duration(15+5*(attempts-5)) * time.Minute
		h.lockoutTimes[clientID] = time.Now().Add(lockoutDuration)
		
		// Log the lockout
		log.Printf("Client %s locked out for %s due to %d failed " + action + "attempts", 
			clientID, lockoutDuration, attempts)
	}
}

// ClearFailedAttempts resets the failed attempts counter for a client
func (h *JsonHandlers) ClearFailedAttempts(clientID string, action string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	delete(h.failedAttempts, clientID)
	delete(h.lockoutTimes, clientID)
	
	// Log successful
	log.Printf("Successful " + action + " for client %s", clientID)
}

// UpdateLastLoginTime records the time of the current login attempt
func (h *JsonHandlers) UpdateLastLoginTime(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.lastLoginTimes[clientID] = time.Now()
}

// LogSuspiciousActivity logs suspicious login activity for monitoring
func (h *JsonHandlers) LogSuspiciousActivity(clientID, username, reason string,action string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	attempts := h.failedAttempts[clientID]
	
	// Log to system log
	log.Printf("Suspicious " + action + " activity: client=%s, username=%s, reason=%s, attempts=%d",
		clientID, username, reason, attempts)
	
	// alerts for high-risk activity
	if attempts >= 3 {

		log.Printf("High-risk " + action + " activity detected for username=%s", username)
	}




	
}




// IncrementSuccessfulRegistrations increments the count of successful registrations for an IP
func (h *JsonHandlers) IncrementSuccessfulRegistrations(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.successfulRegistrations[clientID]++
	count := h.successfulRegistrations[clientID]
	
	log.Printf("Successful registration count for IP %s: %d", clientID, count)
}

// IsRegistrationThresholdReached checks if an IP has reached the threshold for registrations
func (h *JsonHandlers) IsRegistrationThresholdReached(clientID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	count := h.successfulRegistrations[clientID]
	return count >= 3
}

// LockIPForRegistration locks an IP from creating more accounts
func (h *JsonHandlers) LockIPForRegistration(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Set lockout for 24 hours
	lockoutDuration := 24 * time.Hour
	h.lockoutTimes[clientID] = time.Now().Add(lockoutDuration)
	
	log.Printf("IP %s locked from registration for %s due to multiple account creation", 
		clientID, lockoutDuration)
}