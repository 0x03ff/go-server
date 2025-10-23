// internal/store/repositories/users_repository.go
package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

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


func HashPassword(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CompareHashAndData(hash string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
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
		return  fmt.Errorf("failed to hash password")
	}
	user.Password = hashedPassword

	_, err = r.db.Exec(ctx,
		`INSERT INTO users (id, username, password, created_at)
         VALUES (DEFAULT,$1, $2, DEFAULT)`,
		user.Username, user.Password,
	)
	return err
}


func (r *usersRepository) Login(ctx context.Context, username string, password string) (*models.User, error) {
    var user models.User
    err := r.db.QueryRow(ctx,
        `SELECT id, username, password FROM users WHERE username = $1`,
        username,
    ).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    if !CompareHashAndData(user.Password, password) {
        return nil, fmt.Errorf("invalid password")
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
	).Scan( &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}