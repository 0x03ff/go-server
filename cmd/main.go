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
	log.Print("System Logging initialized successfully\n\n")

	// Setup CSV security logging
	csv_filename := "security_events"
	file_number := "_0001"
	csv_logname := csv_filename + file_number
	securityCSV, csvFile, err := utils.SetupSecurityCSV(csv_logname)
	if err != nil {
		log.Fatalf("Failed to initialize security CSV logging: %v", err)
	}
	defer csvFile.Close()

	// Check if Role D testing mode is enabled
	roleDMode := os.Getenv("ROLE_D_MODE") == "true"
	if roleDMode {
		log.Println("⚠️  ROLE D TESTING MODE ENABLED - Authentication disabled for folder downloads")
	} else {
		log.Println("✓ Normal mode - Full authentication enabled")
	}
	log.Println("Server configuration: HTTP:80, HTTPS:443, pprof:8086")


	// set drop_flag to drop the database:
	drop_flag := true

	// set Random_request_address to product different ip-addr:

	random_request_address := true
	if drop_flag {
		// Delete all folders and their contents under /assets/users/
		err = utils.DeleteDirectoryContents("assets/users")
		if err != nil {
			log.Fatal("Failed to delete user assets: ", err)
		}
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
			MaxOpenConns: 50,
			MaxIdleConns: 10,
			MaxIdleTime:  "15m",
		},
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
		JsonHandlers: json_handler.NewHandlers(PDPool, random_request_address, securityCSV),
		RoleDMode:    roleDMode,
	}

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
