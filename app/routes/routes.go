package routes

import (
  "github.com/go-chi/chi"
  "github.com/rbo13/write-it/app"
)

// User sets the user related routes
func User(r chi.Router, handler app.UserHandler) chi.Router {
  r.Get("/", handler.Get)
  r.Get("/{id}", handler.GetByID)
  r.Get("/{id}/posts", handler.GetUserPosts)
  r.Put("/{id}", handler.Update)
  r.Delete("/{id}", handler.Delete)

  // r.Route("/{id}", func(r chi.Router) {
  //  r.Get("/", handler.GetByID)
  //  r.Post("/", handler.Delete)
  //  r.Post("/update", handler.Update)
  // })
  return r
}

// Post sets the post related routes
func Post(r chi.Router, handler app.Handler) chi.Router {

  r.Post("/create", handler.Create)
  r.Get("/", handler.Get)
  r.Get("/{id}", handler.GetByID)
  r.Put("/{id}", handler.Update)
  r.Delete("/{id}", handler.Delete)

  // r.Route("/{id}", func(r chi.Router) {
  //  r.Get("/", handler.GetByID)
  //  r.Post("/", handler.Delete)
  //  r.Post("/update", handler.Update)
  // })
  return r
}
