package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/server/http/middlewares"
	"github.com/ervand7/urlshortener/internal/views/http"
)

// NewRouter creates new chi.Router
func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipMiddleware)

	server := http.Server{
		Storage: storage.GetStorage(),
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/", server.ShortenURL)
		r.Get("/{id:[a-zA-Z]+}", server.GetURL)
		r.Post("/api/shorten", server.APIShortenURL)
		r.Get("/api/user/urls", server.UserURLs)
		r.Post("/api/shorten/batch", server.APIShortenBatch)
		r.Delete("/api/user/urls", server.UserURLsDelete)
		r.Get("/ping", server.PingDB)
	})

	StatsGroup := r.Group(nil)
	StatsGroup.Use(middlewares.NewTrustedNetwork().Handler)
	StatsGroup.Post("/api/internal/stats", server.Stats)

	return r
}
