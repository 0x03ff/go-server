package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// GenerateECCKeyPair returns two strings: the first is the base64-encoded public key, private key.
func GenerateECCKeyPair() (public_key string, private_key string, err error) {
	// Generate a private key for the P-384 curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// The public key is derived from the private key
	publicKey := &privateKey.PublicKey

	// Marshal private key to ASN.1 DER (PKCS#8 format)
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}

	// Encode private key to PEM format
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})

	// Marshal public key to ASN.1 DER (PKIX format)
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	// Encode public key to PEM format
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	// Encode PEM-formatted keys to base64
	base64PublicEncoded := base64.StdEncoding.EncodeToString(pubPEM)
	base64PrivateEncoded := base64.StdEncoding.EncodeToString(privPEM)

	return base64PublicEncoded, base64PrivateEncoded, nil
}

