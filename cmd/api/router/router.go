// router/router.go
package router

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "log"
    "net/http"
    "net/http/pprof"
    "time"
)

type SetupRoutesFunc func(r chi.Router)

func SetupRoutes(setupRoutes SetupRoutesFunc) http.Handler {
    r := chi.NewRouter()

    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)

    r.Use(middleware.Recoverer)
    

    // Set a timeout value on the request context (ctx), that will signal
    // through ctx.Done() that the request has timed out and further
    // processing should be stopped.
    r.Use(middleware.Timeout(60 * time.Second))

    // load static file
    fs := http.FileServer(http.Dir("./assets/"))
    log.Println("Serving files from ./assets/")

    r.Mount("/assets/", http.StripPrefix("/assets/", fs))

    // Mount pprof routes for performance monitoring
    r.Mount("/debug/pprof", http.HandlerFunc(pprof.Index))
    r.Mount("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
    r.Mount("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
    r.Mount("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
    r.Mount("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

    r.Route("/", func(r chi.Router) {
        setupRoutes(r)
    })
    

    return r
}
