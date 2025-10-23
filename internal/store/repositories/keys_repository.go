package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0x03ff/golang/internal/store/models"
)

// Define the keysRepository struct
type keysRepository struct {
	db *pgxpool.Pool
}

// Constructor for keysRepository
func NewKeysRepository(db *pgxpool.Pool) *keysRepository {
	return &keysRepository{db: db}
}

// GetPublicKey retrieves the public key for the specified user.
func (r *keysRepository) GetPublicKey(ctx context.Context, system *models.SystemKey) (string, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, public_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve public key: %w", err)
	}

	return string(system.PublicKey), nil
}

// GetPrivateKey retrieves the private key for the specified user.
func (r *keysRepository) GetPrivateKey(ctx context.Context, system *models.SystemKey) (string, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, private_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %w", err)
	}

	return string(system.PrivateKey), nil
}






