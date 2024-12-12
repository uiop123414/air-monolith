package app

import (
	md "air-monolith/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Route("/api/v1", func(mux chi.Router) {
		mux.Use(md.TimeoutMiddleware(app.Cfg.Timeout))
		mux.Post("/process/sale", app.Sale)
		mux.Post("/process/refund", app.Refund)
	})

	return mux
}
