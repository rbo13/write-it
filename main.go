package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/rbo13/write-it/app/persistence/sql"
	"github.com/rbo13/write-it/app/routes"
	"github.com/rbo13/write-it/app/usecase"
	"github.com/rbo13/write-it/server"
)

const dbName = "writeit"

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

	check(err)

	defer db.Sqlx.Close()

	db.Create(dbName)
	db.Use(dbName)
	db.Migrate()

	// inmemory := inmemory.NewInMemoryPostService()
	// postUsecase := usecase.NewPost(inmemory)

	sqlSrvc := sql.NewSQLService(db.Sqlx)

	postUsecase := usecase.NewPost(sqlSrvc)
	userUsecase := usecase.NewUser(sqlSrvc)

	router.Mount("/api/v1/posts", routes.Routes(router, postUsecase))
	router.Mount("/api/v1/users", routes.Routes(router, userUsecase))
	s := server.New(":1333", router)
	s.StartTLS("./certificates/localhost+2.pem", "./certificates/localhost+2-key.pem")
}

func check(err error) error {
	if err != nil {
		log.Printf("Error occured due to: %v\n\n", err)
		return err
	}

	return nil
}
