package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID uuid.UUID `json:"ID"` // Changed from int64 to uuid.UUID to match PostgreSQL UUID type

	Username string `json:"username"`
	Password string `json:"password"` // Hashed password

	Recover string `json:"recover"` // token to recover key

    ECDH_PublicKey []byte    `json:"ecdh_public_key"` // Serialized public parameters (base64 or hex encoded)

	ECDH_PrivateKey []byte `json:"ecdh_encrypt_private_key"`

	CreatedAt time.Time `json:"created_at" db:"created_at"` // Changed from string to time.Time to properly handle PostgreSQL TIMESTAMPTZ
}
