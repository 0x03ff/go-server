package db

import (
	"context"
	
	"log"
	
	"time"

		
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a native pgx connection pool (no database/sql)
// Perfectly configured for your Docker PostgreSQL container
func NewPostgresPool(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string,drop_flag bool) (*pgxpool.Pool, error) {
	// 1. Parse connection string (handles your Docker credentials)
	poolConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	// 2. Configure pool using pgx's NATIVE settings (critical for Docker)
	poolConfig.MaxConns = int32(maxOpenConns)  // Total max connections
	poolConfig.MinConns = int32(maxIdleConns)  // Minimum warm/idle connections
	poolConfig.MaxConnLifetime = 0             // Let PostgreSQL handle session duration
	poolConfig.HealthCheckPeriod = 1 * time.Minute // Verify connections regularly

	// Parse idle time safely (your Docker needs this!)
	idleDuration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConnIdleTime = idleDuration

	// 3. Create the pool (direct PostgreSQL protocol)
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// 4. Verify connection (with Docker-friendly timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Critical: cleanup if connection fails
		return nil, err
	}

	//5. drop_flag to renew the databae
	if err := db_init(drop_flag, pool); err != nil {
		pool.Close() // Critical: if drop is fail 
		return nil, err
	}


	log.Printf("Connected to PostgreSQL! Pool: %d/%d conns (min/max)", 
		poolConfig.MinConns, poolConfig.MaxConns)
	return pool, nil
}

