package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0x03ff/golang/internal/store/models"
)

// Define the foldersRepository struct
type foldersRepository struct {
	db *pgxpool.Pool
}

// Constructor for foldersRepository
func NewFoldersRepository(db *pgxpool.Pool) *foldersRepository {
	return &foldersRepository{db: db}
}

func (r *foldersRepository) Upload(ctx context.Context, folder *models.Folder) error {
    // Check if folder title already exists
    var exists bool
    err := r.db.QueryRow(ctx,
        `SELECT EXISTS(SELECT 1 FROM folders WHERE title = $1)`,
        folder.Title,
    ).Scan(&exists)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("folder with this name already exists")
    }

    // Prepare the SQL query with all columns
    query := `
        INSERT INTO folders (title, user_id, file_path, secret, encrypt, extension)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at;
    `

    // Execute the query and scan the results
    err = r.db.QueryRow(ctx, query, 
        folder.Title, 
        folder.UserID, 
        folder.FilePath,
        folder.Secret,
        folder.Encrypt,
        folder.Extension).
        Scan(&folder.ID, &folder.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to insert folder record: %w", err)
    }

    return nil
}

func (r *foldersRepository) Search(ctx context.Context, userID uuid.UUID, index int) ([]*models.Folder, error) {
    
    limit := 10
    
    // Calculate offset by aligning to page boundaries
    offset := (index / limit) * limit
        
    query := `
        SELECT id, title, user_id, file_path, encrypt, extension, created_at
        FROM folders
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

    folders := make([]*models.Folder, 0)

    for rows.Next() {
        folder := &models.Folder{}
        err := rows.Scan(
            &folder.ID,
            &folder.Title,
            &folder.UserID,
            &folder.FilePath,
            &folder.Encrypt,
            &folder.Extension,
            &folder.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        folders = append(folders, folder)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return folders, nil
}

func (r *foldersRepository) GetFolderById(ctx context.Context, folder *models.Folder, id int) error {
    query := `
        SELECT id, title, user_id, file_path, secret, encrypt, extension, created_at
        FROM folders
        WHERE id = $1;
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &folder.ID,
        &folder.Title,
        &folder.UserID,
        &folder.FilePath,
        &folder.Secret,
        &folder.Encrypt,
        &folder.Extension,
        &folder.CreatedAt,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return fmt.Errorf("folder with ID %d not found", id)
        }
        return fmt.Errorf("database error: %w", err)
    }
    
    return nil
}