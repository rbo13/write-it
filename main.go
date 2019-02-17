package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	jwtService := jwtservice.New()
	userSQLSrvc := sql.NewUserSQLService(db.Sqlx, jwtService)
	postSQLSrvc := sql.NewPostSQLService(db.Sqlx)

	userUsecase := usecase.NewUser(userSQLSrvc)
	postUsecase := usecase.NewPost(postSQLSrvc)

	router.Post("/register", userUsecase.Create)
	router.Post("/login", userUsecase.Login)

	// Protected routes (API Group)
	router.Group(func(r chi.Router) {
		// Boot up JWT middleware
		r.Use(jwtauth.Verifier(jwtService.TokenAuth))
		r.Use(jwtauth.Authenticator)

		// API GROUP
		// r.Mount("/api/v1/users", routes.User(router, userUsecase))
		// r.Mount("/api/v1/posts", routes.Post(r, postUsecase))

		r.Route("/api", func(rt chi.Router) {
			rt.Mount("/v1/users", routes.User(rt, userUsecase))
			rt.Mount("/v1/posts", routes.Post(r, postUsecase))
		})

		// r.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
		// 	_, claims, _ := jwtauth.FromContext(r.Context())
		// 	w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["email"])))
		// })
	})

	s := server.New(":1333", router)
	go func() {
		s.StartTLS("./certificates/localhost+2.pem", "./certificates/localhost+2-key.pem")
	}()

	gracefulShutdown(s.HTTPServer)

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// s.HTTPServer.Shutdown(ctx)
}

func check(err error) error {
	if err != nil {
		log.Printf("Error occured due to: %v\n\n", err)
		return err
	}

	return nil
}

func gracefulShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
