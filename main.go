package main

import (
	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/rbo13/write-it/app/inmemory"
	"github.com/rbo13/write-it/app/routes"
	"github.com/rbo13/write-it/app/usecase"
	"github.com/rbo13/write-it/server"
)

func main() {
	router := chi.NewRouter()
	// Setup different middlwares here
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		mw.Logger,
		mw.DefaultCompress,
		mw.Recoverer,
		mw.RedirectSlashes,
	)

	inmemory := inmemory.NewInMemoryPostService()
	postUsecase := usecase.NewPost(inmemory)

	router.Mount("/api/v1/posts", routes.Routes(router, postUsecase))
	s := server.New(":1333", router)
	s.Start()
}
