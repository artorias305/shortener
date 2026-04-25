package router

import (
	"shortener/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(urlHandler *handler.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", urlHandler.Shorten)
	})

	r.Get("/{shortCode}", urlHandler.Redirect)

	return r
}
