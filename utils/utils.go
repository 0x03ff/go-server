package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0x03ff/golang/internal/store"
	"github.com/0x03ff/golang/internal/store/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Envelope is a generic type for standardizing API responses
type Envelope map[string]interface{}

// CustomError represents a structured error with status code and details
type CustomError struct {
	Status  int
	Message string
	Details interface{}
}

func (e *CustomError) Error() string {
	return e.Message
}
func MessageToUser(messageToUser string, locationPage string) {
	safeMessage := html.EscapeString(messageToUser)
	fmt.Printf("<script>alert('%s'); window.location.href='../../%s';</script>", safeMessage, locationPage)
}

func SendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Envelope{"error": message})
}

func HashFile(data string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func ValidateInput(description string, data string, lower_limit int, upper_limit int) error {
	if len(data) < lower_limit || len(data) > upper_limit {
		temp := description + " must be between " + strconv.Itoa(lower_limit) + " and " + strconv.Itoa(upper_limit) + " characters"
		return errors.New(temp)
	}

	return nil
}

func GenerateToken(ctx context.Context, userID uuid.UUID, userName string, systemRepo store.SystemKeyRepository) (string, error) {
	const JWT_EXPIRATION = time.Hour * 24

	// Fetch the private key from the system repository

	var systemKey models.SystemKey

	privateKeyPem, err := systemRepo.GetPrivateKey(ctx, &systemKey)
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	// Debugging: Print the retrieved private key
	fmt.Printf("Retrieved Private Key: %s\n", privateKeyPem)

	// Decode the base64-encoded private key
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyPem)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	// Debugging: Print the decoded private key bytes
	fmt.Printf("Decoded Private Key Bytes: %x\n", privateKeyBytes)

	// Decode the PEM block
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}

	// Debugging: Print the PEM block type and bytes
	fmt.Printf("PEM Block Type: %s\n", block.Type)
	fmt.Printf("PEM Block Bytes: %x\n", block.Bytes)

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

    publicKeyPem, err := systemRepo.GetPublicKey(context.Background(), systemKey)

	
    if err != nil {
        return nil, fmt.Errorf("failed to get public key: %w", err)
    }


    // Decode the base64-encoded public key
    publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyPem)
    if err != nil {
        return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
    }


    // Decode the PEM block
    block, rest := pem.Decode(publicKeyBytes)
    if block == nil || len(rest) > 0 {
        return nil, errors.New("failed to decode public key PEM")
    }

    log.Printf("Decoded PEM block type: %s", block.Type)

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









