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

// GetECDSAPublicKey retrieves the ECDSA public key for the specified user.
func (r *keysRepository) GetECDSAPublicKey(ctx context.Context, system *models.SystemKey) ([]byte, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, public_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve system public key: %w", err)
	}

	return system.PublicKey, nil
}

// GetECDSAPrivateKey retrieves the ECDSA private key for the server.
func (r *keysRepository) GetECDSAPrivateKey(ctx context.Context, system *models.SystemKey) ([]byte, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, private_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve system private key: %w", err)
	}

	return system.PrivateKey, nil
}








// GetECDHPublicKey retrieves the ECDH public key for the specified user.
func (r *keysRepository) GetECDHPublicKey(ctx context.Context, system *models.SystemKey) ([]byte, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, ecdh_public_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.ECDH_PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve system public key: %w", err)
	}

	return system.ECDH_PublicKey, nil
}

// GetECDHPrivateKey retrieves the ECDH private key for the server.
func (r *keysRepository) GetECDHPrivateKey(ctx context.Context, system *models.SystemKey) ([]byte, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, ecdh_private_key 
		FROM system_keys 
		LIMIT 1
	`)

	err := row.Scan(&system.ID, &system.ECDH_PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve system private key: %w", err)
	}

	return system.ECDH_PrivateKey, nil
}

