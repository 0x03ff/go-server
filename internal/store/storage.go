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
		Login(ctx context.Context, username string, password string,recover string) (*models.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
		GetUserECDHPrivateKey(ctx context.Context,user *models.User) ([]byte, error)
		GetUserECDHPublicKey(ctx context.Context,user *models.User) ([]byte, error)

	}

type SystemKeyRepository interface {
		GetECDSAPublicKey(ctx context.Context, system *models.SystemKey) ([]byte,error)
		GetECDSAPrivateKey(ctx context.Context, system *models.SystemKey) ([]byte,error)
		GetECDHPublicKey(ctx context.Context, system *models.SystemKey) ([]byte, error)
		GetECDHPrivateKey(ctx context.Context, system *models.SystemKey) ([]byte, error)
	}
type FileRepository interface {
	 Upload(ctx context.Context, user *models.File) (err error)
	 Search(ctx context.Context, userID uuid.UUID, index int) ([]*models.File, error) 
	 
	 GetFileById(ctx context.Context, file *models.File, id int) error
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