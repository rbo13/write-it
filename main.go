package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/rbo13/write-it/app/persistence/inmemory"
	"github.com/rbo13/write-it/app/persistence/sql"
	"github.com/rbo13/write-it/app/routes"
	"github.com/rbo13/write-it/app/usecase"
	"github.com/rbo13/write-it/server"
)

// const dbName = "write-it"

const dsn = "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4"

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

	db, err := sql.New(dsn)

	if err != nil {
		log.Printf("Error New: %v", err)
		return
	}

	log.Println(db)

	inmemory := inmemory.NewInMemoryPostService()
	postUsecase := usecase.NewPost(inmemory)

	router.Mount("/api/v1/posts", routes.Routes(router, postUsecase))
	s := server.New(":1333", router)
	s.Start()
}
