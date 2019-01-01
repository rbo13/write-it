package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/rbo13/write-it/app/jwtservice"
	"github.com/rbo13/write-it/app/persistence/sql"
	"github.com/rbo13/write-it/app/routes"
	"github.com/rbo13/write-it/app/usecase"
	"github.com/rbo13/write-it/server"
)

const (
	dbName = "writeit"
	dsn    = "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4"
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

	db, err := sql.New(dsn)

	check(err)

	defer db.Sqlx.Close()

	db.Create(dbName)
	db.Use(dbName)
	db.Migrate()

	// inmemory := inmemory.NewInMemoryPostService()
	// postUsecase := usecase.NewPost(inmemory)

	jwtService := jwtservice.New()
	sqlSrvc := sql.NewSQLService(db.Sqlx, jwtService)

	postUsecase := usecase.NewPost(sqlSrvc)
	userUsecase := usecase.NewUser(sqlSrvc)

	// Protected routes (API Group)
	router.Group(func(r chi.Router) {
		// TODO :: Seek, verify and validate JWT tokens
		// custom jwt verifier middleware
		r.Use(jwtauth.Verifier(jwtService.TokenAuth))

		// TODO :: Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		// custom jwt middleware authenticator
		r.Use(jwtauth.Authenticator)

		// API GROUP
		r.Mount("/api/v1/users", routes.User(router, userUsecase))
		r.Mount("/api/v1/posts", routes.Post(router, postUsecase))

		r.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["email"])))
		})
	})

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
