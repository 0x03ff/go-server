package models

import (
	"time"
	"github.com/google/uuid"
)

type File struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    UserID    uuid.UUID `json:"user_id"`
    FilePath  string    `json:"file_path"`
    Extension string    `json:"extension"`// Stores file extension (.txt, .pdf, etc.)
    CreatedAt time.Time `json:"created_at"`
}

type FilesResponse struct {
    Files    []*File `json:"files"`
    Index    int     `json:"index"`
    Page     int     `json:"page"`
    HasPrev  bool    `json:"has_prev"`
    HasNext  bool    `json:"has_next"`
}