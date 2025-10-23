package store

import (
	"context"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


type UsersRepository interface {
		Create(context.Context, *models.User) error
		Login(ctx context.Context, username string, password string) (*models.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	}

type SystemKeyRepository interface {
		GetPublicKey(ctx context.Context, user *models.SystemKey) (string,error)
		GetPrivateKey(ctx context.Context, user *models.SystemKey) (string,error)
		
	}
type FileRepository interface {
	 Upload(ctx context.Context, user *models.File) (err error)
}

type Storage struct {
	
	Users UsersRepository
	System SystemKeyRepository
	Files FileRepository
} 

func NewStorage(db *pgxpool.Pool) Storage{
	return Storage{
		
		Users: repositories.NewUsersRepository(db),
		System: repositories.NewKeysRepository(db),
		Files : repositories.NewFilesRepository(db),
	}
}