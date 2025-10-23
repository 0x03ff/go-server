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
	Sysconfig Config
	Store     store.Storage
	Tlsconfig *tls.Config
	Cert_path string
	Key_path  string
	HtmlHandlers  *html_handler.HtmlHandlers
	JsonHandlers   *json_handler.JsonHandlers
	
}

func (app *Application) GetHtmlHandlers() *html_handler.HtmlHandlers {
    return app.HtmlHandlers
}

func (app *Application) GetJsonHandlers() *json_handler.JsonHandlers {
    return app.JsonHandlers
}





func (app *Application) Mount() http.Handler {
    setupFunc := func(r chi.Router) {

        html_handler.SetupHtmlRoutes(r, app)

		json_handler.SetupJsonRoutes(r, app)
    }
    return router.SetupRoutes(setupFunc)
}

func (app *Application) Run(mux http.Handler) error {
    srv := &http.Server{
        Addr:         app.Sysconfig.ADDR,
        Handler:      mux,
        WriteTimeout: time.Second * 30,
        ReadTimeout:  time.Second * 10,
        IdleTimeout:  time.Minute,
        TLSConfig:    app.Tlsconfig,
    }

    log.Printf("Server has started at %s",
        app.Sysconfig.ADDR)

    return srv.ListenAndServeTLS(app.Cert_path, app.Key_path)
}
