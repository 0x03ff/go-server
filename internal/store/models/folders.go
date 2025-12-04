package models

import (
	"time"

	"github.com/google/uuid"
)

// For more file / contain in folder

type Folder struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	UserID     uuid.UUID `json:"user_id"`
	FilePath   string    `json:"file_path"`
	Secret     []byte    `json:"-"`         // store the encrypted AES key (for RSA) or AES key directly
	PrivateKey []byte    `json:"-"`         // store the RSA private key (only for RSA encryption)
	Encrypt    string    `json:"encrypt_method"` // store the encrypt method (non-encrypted, aes, rsa-2048, rsa-4096, etc.)
	Extension  string    `json:"extension"` // Stores file extension (.txt, .pdf, etc.)
	CreatedAt  time.Time `json:"created_at"`
}

type FoldersResponse struct {
    Folders  []*Folder `json:"folders"`
    Index    int       `json:"index"`
    Page     int       `json:"page"`
    HasPrev  bool      `json:"has_prev"`
    HasNext  bool      `json:"has_next"`
}
