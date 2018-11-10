package routes

import (
	"github.com/go-chi/chi"
	"github.com/rbo13/write-it/app"
)

// Routes ...
func Routes(r chi.Router, handler app.Handler) chi.Router {

	r.Post("/create", handler.Create)
	r.Get("/", handler.Get)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Post("/", handler.Delete)
		r.Post("/update", handler.Update)
	})
	return r
}
