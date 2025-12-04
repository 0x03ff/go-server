package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	_ "net/http/pprof" // Import this for pprof endpoints

	"github.com/0x03ff/golang/cmd/api/config"
	"github.com/0x03ff/golang/cmd/api/router/html_handler"
	"github.com/0x03ff/golang/cmd/api/router/json_handler"
	"github.com/0x03ff/golang/internal/db"
	"github.com/0x03ff/golang/internal/store"
	"github.com/0x03ff/golang/utils"
)

func main() {

	_, logFile, err := utils.SetupLogging()
	if err != nil {
		// Fallback to stderr since logging isn't set up yet
		log.Fatalf("Failed to initialize logging: %v", err)
	}
	defer logFile.Close()
	// ALL subsequent logs go to BOTH terminal and logs/system.log
	log.Print("Logging initialized successfully\n\n")

	// Check if Role D testing mode is enabled
	roleDMode := os.Getenv("ROLE_D_MODE") == "true"
	if roleDMode {
		log.Println("⚠️  ROLE D TESTING MODE ENABLED - Authentication disabled for folder downloads")
		log.Println("⚠️  HTTP and HTTPS servers will run on ports 80 and 443")
	} else {
		log.Println("✓ Normal mode - Full authentication enabled, HTTPS only on port 443")
	}

	//TLS config
	tlsCfg := &config.TlsConfig{
		CertFile:   "internal/certs/go_cert.pem",
		KeyFile:    "internal/certs/go_key.pem",
		MinVersion: tls.VersionTLS12,
	}

	tls_config, err := tlsCfg.NewTLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Config{
		ADDR: "0.0.0.0:443",
		DB: config.DBConfig{
			DB_addr:      "postgres://comp4334:secret@localhost:5432/go_server?sslmode=disable",
			MaxOpenConns: 25,
			MaxIdleConns: 5,
			MaxIdleTime:  "15m",
		},
	}

	// set drop_flag to drop the database:
	drop_flag := true

	if drop_flag {
		// Delete all folders and their contents under /assets/users/
		err = utils.DeleteDirectoryContents("assets/users")
		if err != nil {
			log.Fatal("Failed to delete user assets: ", err)
		}
	}

	PDPool, err := db.NewPostgresPool(cfg.DB.DB_addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime, drop_flag)
	if err != nil {
		log.Fatal("DB connection failed: ", err)
	}
	defer PDPool.Close()

	store := store.NewStorage(PDPool)

	// Verify database connection
	if err := PDPool.Ping(context.Background()); err != nil {
		log.Fatal("Could not ping database: ", err)
	}
	log.Println("Database connection established successfully")

	app := &config.Application{
		Sysconfig:    cfg,
		Store:        store,
		Tlsconfig:    tls_config,
		Cert_path:    "internal/certs/go_cert.pem",
		Key_path:     "internal/certs/go_key.pem",
		HtmlHandlers: html_handler.NewHandlers(PDPool),
		JsonHandlers: json_handler.NewHandlers(PDPool),
		RoleDMode:    roleDMode,
	}

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
