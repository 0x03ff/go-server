package models

import (
	"time"

	"github.com/google/uuid"
)



type File struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	UserID uuid.UUID `json:"user_id"` // Data owner ID

	Share []byte `json:"parameter"` 
	File string `json:"file"` // store the file path in db
	CreatedAt time.Time `json:"created_at"`
}
