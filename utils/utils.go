package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/0x03ff/golang/internal/store"
	"github.com/0x03ff/golang/internal/store/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

//1. user function

// Envelope is a generic type for standardizing API responses

// CustomError represents a structured error with status code and details
type CustomError struct {
	Status  int
	Message string
	Details interface{}
}

func (e *CustomError) Error() string {
	return e.Message
}

func SendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})

}

func SendJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   message,
		"code":    code,
		"details": message,
	})
}



func SanitizeHTML(input string) string {
	// Sanitize the input using bluemonday
	policy := bluemonday.UGCPolicy()
	clean := policy.Sanitize(input)
	return clean
}

func ValidateUserInput(description string, data string, lower_limit int, upper_limit int) error {

	// Sanitize the input to prevent XSS
	sanitizedData := SanitizeHTML(data)

	if len(sanitizedData) < lower_limit || len(sanitizedData) > upper_limit {
		temp := description + " must be between " + strconv.Itoa(lower_limit) + " and " + strconv.Itoa(upper_limit) + " characters"
		return errors.New(temp)
	}

	return nil
}

// 2. user JWT
func GenerateToken(ctx context.Context, userID uuid.UUID, userName string, systemRepo store.SystemKeyRepository) (string, error) {
	const JWT_EXPIRATION = time.Hour * 24

	// Fetch the private key from the system repository

	var systemKey models.SystemKey

	privateKeyPem, err := systemRepo.GetECDSAPrivateKey(ctx, &systemKey)
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(privateKeyPem)
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}

	// Parse the private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Add detailed error message
		return "", fmt.Errorf("failed to parse private key: %v, raw bytes: %x", err, block.Bytes)
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"exp":       time.Now().Add(JWT_EXPIRATION).Unix(),
	})

	// Sign and get the complete encoded token as a string using the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		// Add detailed error message
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, systemRepo store.SystemKeyRepository) (*jwt.Token, error) {

	// Fetch the public key from the system repository
	systemKey := &models.SystemKey{}

	publicKeyPem, err := systemRepo.GetECDSAPublicKey(context.Background(), systemKey)

	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Decode the PEM block
	block, rest := pem.Decode(publicKeyPem)
	if block == nil || len(rest) > 0 {
		return nil, errors.New("failed to decode public key PEM")
	}

	// Parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Ensure the public key is of the correct type
	pubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("Public key is not of type *ecdsa.PublicKey")
		return nil, errors.New("public key is not of type *ecdsa.PublicKey")
	}

	// Parse the token
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		log.Printf("Invalid token")
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// 3. system function

func HashFile(data string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func DeleteDirectoryContents(dirPath string) error {
	// Read the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory does not exist, nothing to delete
		}
		return err
	}

	// Remove each entry in the directory
	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			// Recursively delete subdirectories
			err = os.RemoveAll(entryPath)
			if err != nil {
				return err
			}
		} else {
			// Delete files
			err = os.Remove(entryPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GenerateCSRFToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to generate CSRF token: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(b)
}

func VerifyCSRFtoken(w http.ResponseWriter, r *http.Request) error {
	// 1. Try to get token from header first (for AJAX requests)
	csrfToken := r.Header.Get("X-CSRF-Token")

	// 2. If not in header, check form data (for traditional form submissions)
	if csrfToken == "" {
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Parse multipart form with reasonable size limit
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				return fmt.Errorf("failed to parse multipart form: %w", err)
			}
			csrfToken = r.FormValue("csrf_token")
		} else {
			// For URL-encoded form data
			if err := r.ParseForm(); err != nil {
				return fmt.Errorf("failed to parse form: %w", err)
			}
			csrfToken = r.FormValue("csrf_token")
		}
	}

	// 3. Validate we have a token
	if csrfToken == "" {
		SendError(w, http.StatusForbidden, "Missing CSRF token")
		return errors.New("missing csrf token")
	}

	// 4. Get token from cookie
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil {
		SendError(w, http.StatusForbidden, "Invalid CSRF token")
		return fmt.Errorf("csrf cookie error: %w", err)
	}

	// 5. Validate token match
	if csrfCookie.Value != csrfToken {
		SendError(w, http.StatusForbidden, "Invalid CSRF token")
		return errors.New("csrf token mismatch")
	}

	return nil
}

func SetupLogging() (*log.Logger, *os.File, error) {
	// Create logs directory if missing
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err
	}

	// Open log file with append mode
	logFile, err := os.OpenFile("./logs/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Create multi-writer (terminal + log file)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure standard logger to use multi-writer
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Optional: add timestamps and file info

	// Return the standard logger, log file, and nil error
	return log.Default(), logFile, nil
}

