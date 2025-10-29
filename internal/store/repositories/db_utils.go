package repositories

import (

	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CompareHashAndData(hash string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// GenerateECDSAKeyPair returns two strings: the first is the public key, then private key.

// The key pair with ECDSA Usage eg JWT

func GenerateECDSAKeyPair() ([]byte, []byte, error) {
    // Generate a private key for the P-384 curve
    privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
    }

    // The public key is derived from the private key
    publicKey := &privateKey.PublicKey

    privBytes, err := x509.MarshalECPrivateKey(privateKey)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
    }

    // Encode private key to PEM format (as raw bytes)
    privPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "EC PRIVATE KEY",
        Bytes: privBytes,
    })

    // Marshal public key to ASN.1 DER (PKIX format)
    pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
    }

    // Encode public key to PEM format (as raw bytes)
    pubPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubBytes,
    })

    // Return raw PEM bytes (not base64 encoded)
    return pubPEM, privPEM, nil
}


func GenerateECDHKeyPair() ([]byte, []byte, error) {
    curve := ecdh.P384()
    
    // Generate a proper ECDH key pair
    privateKey, err := curve.GenerateKey(rand.Reader)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate ECDH private key: %w", err)
    }

    // Get the public key
    publicKey := privateKey.PublicKey()

    // Marshal private key to PKCS#8 format
    privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
    }

    // Encode private key to PEM format
    privPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PRIVATE KEY", // Correct for PKCS#8
        Bytes: privBytes,
    })

    // Marshal public key to PKIX format
    pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
    }
    
    // Encode public key to PEM format 
    pubPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubBytes,
    })

    // Return raw PEM bytes
    return pubPEM, privPEM, nil
}

func GenerateHash256(input string) []byte {
    normalized := strings.TrimSpace(input)
    
    hasher := sha256.New()
    hasher.Write([]byte(normalized))
    return hasher.Sum(nil)
}


