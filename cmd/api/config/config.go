package config

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/0x03ff/golang/cmd/api/router"
	"github.com/0x03ff/golang/cmd/api/router/html_handler"
	"github.com/0x03ff/golang/cmd/api/router/json_handler"
	"github.com/0x03ff/golang/internal/store"
	"github.com/go-chi/chi/v5"
)

//for dependencies

type DBConfig struct {
	DB_addr      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type TlsConfig struct {
	CertFile   string
	KeyFile    string
	MinVersion uint16
	TLSConfig  *tls.Config
}

func (t *TlsConfig) NewTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}

	t.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   t.MinVersion,
	}

	return t.TLSConfig, nil
}

// configuration
type Config struct {
	ADDR string
	DB   DBConfig
}

type Application struct {
	Sysconfig    Config
	Store        store.Storage
	Tlsconfig    *tls.Config
	Cert_path    string
	Key_path     string
	HtmlHandlers *html_handler.HtmlHandlers
	JsonHandlers *json_handler.JsonHandlers
	RoleDMode    bool // Enable Role D testing mode (disable auth for folder downloads)
}

func (app *Application) GetHtmlHandlers() *html_handler.HtmlHandlers {
	return app.HtmlHandlers
}

func (app *Application) GetJsonHandlers() *json_handler.JsonHandlers {
	return app.JsonHandlers
}

func (app *Application) Mount() http.Handler {
	setupFunc := func(r chi.Router) {

		json_handler.SetupJsonRoutes(r, app)

		html_handler.SetupHtmlRoutes(r, app)

	}

	return router.SetupRoutes(setupFunc, app.RoleDMode)
}
func (app *Application) Run(mux http.Handler) error {
	// 1. Start pprof HTTP server (port 8086) in goroutine
	go func() {
		log.Printf("pprof server started on :8086")
		if err := http.ListenAndServe("0.0.0.0:8086", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	// HTTPS server (port 443)
	httpsSrv := &http.Server{
		Addr:         "0.0.0.0:443",
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
		TLSConfig:    app.Tlsconfig,
	}

	// HTTP server (port 80) - always start in normal mode, conditional in Role D mode
	httpSrv := &http.Server{
		Addr:         "0.0.0.0:80",
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Println("HTTP server started on 0.0.0.0:80")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start HTTPS server (blocks main thread)
	log.Println("HTTPS server started on 0.0.0.0:443")
	return httpsSrv.ListenAndServeTLS(app.Cert_path, app.Key_path)
}
