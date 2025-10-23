package repositories

import (
	"context"


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
    // Prepare the SQL query
    query := `
        INSERT INTO files (title, user_id, share_parameter, file_path)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at;
    `

    // Execute the query and scan the results
    err := r.db.QueryRow(ctx, query, file.Title, file.UserID, file.Share, file.File).
        Scan(&file.ID, &file.CreatedAt)

    if err != nil {
        return err
    }

    return nil
}
