package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID                uuid.UUID `json:"ID"` // Changed from int64 to uuid.UUID to match PostgreSQL UUID type
    
    Username          string    `json:"username"`
    Password          string    `json:"password"` // Hashed password


    CreatedAt         time.Time `json:"created_at" db:"created_at"` // Changed from string to time.Time to properly handle PostgreSQL TIMESTAMPTZ
}



