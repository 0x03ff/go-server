// internal/store/repositories/users_repository.go
package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0x03ff/golang/internal/store/models"
)

// 1. Define the struct
type usersRepository struct {
	db *pgxpool.Pool
}

// 2. Constructor function
func NewUsersRepository(db *pgxpool.Pool) *usersRepository {
	return &usersRepository{db: db}
}

// 3. Implement methods
func (r *usersRepository) Create(ctx context.Context, user *models.User) error {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		user.Username,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user name already exists")
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	user.Password = hashedPassword

	hashedRecover, err := HashPassword(user.Recover)
	if err != nil {
		return fmt.Errorf("failed to hash Recover token")
	}

	publicKey, privateKey, err := GenerateECDHKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate user ECC key pair: %v", err)
	}

	println("user.Recover:    ", string(user.Recover))


	if err != nil {
		log.Fatalf("Failed to encrypt user ECC key pair: %v", err)
	}

	// user recover is under b-crypt now
	user.Recover = hashedRecover

	_, err = r.db.Exec(ctx,
		`INSERT INTO users (id, username, password,recover,ecdh_public_key,ecdh_encrypt_private_key, created_at)
         VALUES (DEFAULT,$1, $2, $3,$4, $5, DEFAULT)`,
		user.Username, user.Password, user.Recover, publicKey, privateKey,
	)
	return err
}

func (r *usersRepository) Login(ctx context.Context, username string, password string, recover string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, password,recover FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Recover)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if !CompareHashAndData(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	if !CompareHashAndData(user.Recover, recover) {
		return nil, fmt.Errorf("invalid recover key")
	}

	return &user, nil
}

// Get a user by ID
func (r *usersRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}

// Get a user by username
func (r *usersRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username FROM users WHERE username = $1`,
		username,
	).Scan(&user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}

// GetUserECDHPrivateKey retrieves the ECDH private key for the specified user.
func (r *usersRepository) GetUserECDHPrivateKey(ctx context.Context, user *models.User) ([]byte, error) {
	var ecdh_privateKey []byte
	err := r.db.QueryRow(ctx, `
        SELECT ecdh_encrypt_private_key 
        FROM users 
        WHERE id = $1
    `, user.ID).Scan(&ecdh_privateKey)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	return ecdh_privateKey, nil
}

// GetUserECDHPublicKey retrieves the ECDH public key for the specified user.
func (r *usersRepository) GetUserECDHPublicKey(ctx context.Context, user *models.User) ([]byte, error) {
	var ecdh_publicKey []byte
	err := r.db.QueryRow(ctx, `
        SELECT ecdh_public_key 
        FROM users 
        WHERE id = $1
    `, user.ID).Scan(&ecdh_publicKey)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	return ecdh_publicKey, nil
}
