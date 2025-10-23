package models

import (
	"github.com/google/uuid"
)

type SystemKey struct {
	ID        uuid.UUID `json:"id"`                // Public Parameters (PP) - stored in plaintext (it's public)
	PublicKey []byte    `json:"public_key"` // Serialized public parameters (base64 or hex encoded)

	PrivateKey []byte `json:"-"`

}
