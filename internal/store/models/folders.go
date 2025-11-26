package models

import (
	"time"

	"github.com/google/uuid"
)

// For more file / contain in folder

type Folder struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UserID    uuid.UUID `json:"user_id"`
	FilePath  string    `json:"file_path"`
	Secret    []byte    `json:"-"`         // store the encrpyt AES key
	Extension string    `json:"extension"` // Stores file extension (.txt, .pdf, etc.)
	CreatedAt time.Time `json:"created_at"`
}
