package db

import (
	"context"
	"fmt"
	"log"

	"github.com/0x03ff/golang/internal/store"
	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func db_init(dropFlag bool, pool *pgxpool.Pool) error {
	ctx := context.Background()

	var err error

	if dropFlag {
		// Drop the tables if they exist
		tables := []string{"users", "system_keys", "files"}
		for _, table := range tables {
			_, err = pool.Exec(ctx, "DROP TABLE IF EXISTS "+table)
			if err != nil {
				return err
			}

		}
		fmt.Println("Dropped tables")
	}

	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`)
	if err != nil {
		return err
	}

	fmt.Println("Created extension")

	// Create the users table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	fmt.Println("Created users table")

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS system_keys (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			public_key BYTEA NOT NULL,
			private_key BYTEA NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	fmt.Println("Created system_keys table")

	_, err = pool.Exec(ctx, `
		CREATE TABLE files (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title VARCHAR(255) NOT NULL,
			user_id UUID NOT NULL,
			share_parameter BYTEA NOT NULL,
			file_path TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

`)
	if err != nil {
		return err
	}

	fmt.Println("Created files table")

	if dropFlag {
		// Create a new storage instance
		storage := store.NewStorage(pool)

		// Create default users
		ctx := context.Background()

		// Default users to be inserted into the database
		defaultUsers := []struct {
			Username string
			Password string
		}{
			{"testtest", "testtest"},
			{"admin", "admin"},
		}

		for _, user := range defaultUsers {
			newUser := &models.User{
				Username: user.Username,
				Password: user.Password,
			}

			err := storage.Users.Create(ctx, newUser)
			if err != nil {
				log.Fatalf("Failed to insert default user %s: %v", user.Username, err)
			}
		}

		fmt.Println("Inserted default users")

		publicKey, privateKey, err := utils.GenerateECCKeyPair()
		if err != nil {
			log.Fatalf("Failed to generate ECC key pair: %v", err)
		}

		// Insert the keys into the system_keys table
		_, err = pool.Exec(ctx, `
			INSERT INTO system_keys (public_key, private_key)
			VALUES ($1, $2)
		`, []byte(publicKey), []byte(privateKey))
		if err != nil {
			log.Fatalf("Failed to insert keys into system_keys table: %v", err)
		}

		fmt.Println("Inserted keys into system_keys table")
		

	}


	return nil
}
