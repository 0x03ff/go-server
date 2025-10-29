package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0x03ff/golang/internal/store/models"
)

// Define the keysRepository struct
type filesRepository struct {
	db *pgxpool.Pool
}

// Constructor for keysRepository
func NewFilesRepository(db *pgxpool.Pool) *filesRepository {
	return &filesRepository{db: db}
}
func (r *filesRepository) Upload(ctx context.Context, file *models.File) error {
    // Check if file title already exists (including extension)
    var exists bool
    err := r.db.QueryRow(ctx,
        `SELECT EXISTS(SELECT 1 FROM files WHERE title = $1)`,
        file.Title,
    ).Scan(&exists)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("file with this name already exists")
    }

    // Prepare the SQL query with the extension column
    query := `
        INSERT INTO files (title, user_id, file_path, extension)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at;
    `

    // Execute the query and scan the results
    err = r.db.QueryRow(ctx, query, 
        file.Title, 
        file.UserID, 
        file.FilePath,
        file.Extension).  // This is the critical addition
        Scan(&file.ID, &file.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to insert file record: %w", err)
    }

    return nil
}



func (r *filesRepository) Search(ctx context.Context, userID uuid.UUID, index int) ([]*models.File, error) {
    
    limit := 10
    
     // Calculate offset by aligning to page boundaries
    // This ensures that any index within a page range (0-9, 10-19, etc.)
    // returns the same page of results
    offset := (index / limit) * limit
        
    query := `
        SELECT id, title, user_id, file_path, created_at
        FROM files
        WHERE user_id = $1
        ORDER BY created_at DESC
        OFFSET $2
        LIMIT $3;
    `

    rows, err := r.db.Query(ctx, query, userID, offset, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    files := make([]*models.File, 0)

    for rows.Next() {
        file := &models.File{}
        err := rows.Scan(
            &file.ID,
            &file.Title,
            &file.UserID,
            &file.FilePath,
            &file.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        files = append(files, file)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return files, nil
}


func (r *filesRepository) GetFileById(ctx context.Context, file *models.File, id int) error {
    query := `
        SELECT id, title, user_id, file_path, created_at
        FROM files
        WHERE id = $1;
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &file.ID,
        &file.Title,
        &file.UserID,
        &file.FilePath,
        &file.CreatedAt,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return fmt.Errorf("file with ID %d not found", id)
        }
        return fmt.Errorf("database error: %w", err)
    }
    
    return nil
}