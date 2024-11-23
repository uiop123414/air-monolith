package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(mux chi.Router) {
		mux.Use(app.TimeoutMiddleware(app.cfg.timeout))
		mux.Post("/process/sale", app.Sale)
		mux.Post("/process/refund", app.Refund)
	})

	return mux
}
