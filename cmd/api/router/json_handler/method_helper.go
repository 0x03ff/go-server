package json_handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func generateRandomIP() string {

	
	// 1. Legitimate user with around 10 devices (uncomment for normal testing)
	//return fmt.Sprintf("192.168.0.%d", rand.Intn(10)+2)

	// 2. Attacker with around 50 devices
	// return fmt.Sprintf("192.168.0.%d", rand.Intn(50)+2)

	// 3. Attacker with around 100 devices
	// return fmt.Sprintf("192.168.0.%d", rand.Intn(100)+2)

	// 4. Attacker with around 200 devices
	// return fmt.Sprintf("192.168.0.%d", rand.Intn(200)+2)

	// 5. Attacker with around 300 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 6. Attacker with around 500 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 7. Attacker with around 700 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 8. Attacker with around 1000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 9. Attacker with around 2000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 10. Attacker with around 5000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 11. Attacker with around 10000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 12. Attacker with around 20000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 13. Attacker with around 50000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// 14. Attacker with around 100000 devices
	// return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))

	// Extreme case for attack testing (100,000+ distinct IPs)
	// Uses RFC 5737 benchmarking range (198.18.0.0/15) - safe for testing
	return fmt.Sprintf("198.%d.%d.%d", 18 + rand.Intn(2), rand.Intn(256), rand.Intn(256))
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

	if h.random_address {

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
		log.Printf("Applying %d second delay for client "+action+" %s (attempt %d)", delay, clientID, attempts)

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
    log.Printf("Failed "+action+" attempt %d for client %s", attempts, clientID)

    // Apply lockout after 5 failed attempts
    if attempts >= 5 {
        lockoutDuration := time.Duration(15+5*(attempts-5)) * time.Minute
        h.lockoutTimes[clientID] = time.Now().Add(lockoutDuration)

        // Log the lockout with CORRECT format
        log.Printf("Client %s locked out for %s due to %d failed %s attempts",
            clientID, lockoutDuration, attempts, action)
        
        // CSV LOGGING FOR LOGIN LOCKOUT - USE CONSISTENT FORMAT
        h.logSecurityEvent(
            "account_lockout",  // Use consistent event type
            clientID,
            "",
            fmt.Sprintf("Locked out for %s due to %d failed %s attempts", 
                        lockoutDuration, attempts, action),
            "security",         // Action should be "security"
            int(lockoutDuration.Seconds()),  // Store duration in seconds
        )
    }
}


// ClearFailedAttempts resets the failed attempts counter for a client
func (h *JsonHandlers) ClearFailedAttempts(clientID string, username string, action string) {
    h.mu.Lock()
    defer h.mu.Unlock()

    delete(h.failedAttempts, clientID)
    delete(h.lockoutTimes, clientID)

    // Log successful
    log.Printf("Successful "+action+" for client %s", clientID)

    // Record in CSV with proper event type
    eventType := "login_success"
    if action == "register" {
        eventType = "registration_success"
    }
    
    h.logSecurityEvent(
        eventType,
        clientID,
        username,
        "Successful "+action,
        action,
        1,
    )
}
// UpdateLastLoginTime records the time of the current login attempt
func (h *JsonHandlers) UpdateLastLoginTime(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.lastLoginTimes[clientID] = time.Now()
}

// LogSuspiciousActivity logs suspicious login activity for monitoring
func (h *JsonHandlers) LogSuspiciousActivity(clientID, username, reason string, action string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	attempts := h.failedAttempts[clientID]

	// Log to system log
	log.Printf("Suspicious "+action+" activity: client=%s, username=%s, reason=%s, attempts=%d",
		clientID, username, reason, attempts)

	// CSV LOGGING FOR SUSPICIOUS ACTIVITY
	h.logSecurityEvent("suspicious_activity", clientID, username, reason, action, attempts)

	// alerts for high-risk activity
	if attempts >= 3 {

		log.Printf("High-risk "+action+" activity detected for username=%s", username)
	}

}

// IncrementSuccessfulRegistrations increments the count of successful registrations for an IP
func (h *JsonHandlers) IncrementSuccessfulRegistrations(clientID string, username string) {
	h.mu.Lock()
	h.successfulRegistrations[clientID]++
	count := h.successfulRegistrations[clientID]
	h.mu.Unlock()

	log.Printf("Successful registration count for IP %s: %d", clientID, count)

	// CSV LOGGING FOR REGISTRATION ATTEMPTS
	h.logSecurityEvent("registration_attempt", clientID, username, "new account creation", "registration", count)
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
	count := h.successfulRegistrations[clientID]
	lockoutDuration := 24 * time.Hour
	h.lockoutTimes[clientID] = time.Now().Add(lockoutDuration)
	h.mu.Unlock()

	log.Printf("IP %s locked from registration for %s due to multiple account creation",
		clientID, lockoutDuration)

	// CSV LOGGING FOR REGISTRATION LOCK
	h.logSecurityEvent("registration_lock", clientID, "", "multiple account creation", "registration", count)
}

func (h *JsonHandlers) logSecurityEvent(eventType, clientID, username, reason, action string, attempts int) {
    h.csvMu.Lock()
    defer h.csvMu.Unlock()

    timestamp := time.Now().Format("2006-01-02T15:04:05+08:00")
    record := []string{
        timestamp,
        eventType,
        clientID,
        username,
        reason,
        strconv.Itoa(attempts),
        action,
    }
    
    if err := h.csvWriter.Write(record); err != nil {
        log.Printf("Error writing to security CSV: %v", err)
    }
    h.csvWriter.Flush()
}

