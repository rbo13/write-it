package sql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rbo13/write-it/app"
)

var (
	errEmpty       = errors.New("error: Post is required")
	errNotInserted = errors.New("error: Not inserted")
)

// PostService implements the app.UserService
type PostService interface {
	app.PostService
}

// Post implements the PostService interface
type Post struct {
	DB       *sqlx.DB
	PostSrvc *app.Post
}

// NewPostSQLService returns the interface that implements the app.PostService
func NewPostSQLService(db *sqlx.DB) PostService {
	return &Post{
		DB:       db,
		PostSrvc: new(app.Post),
	}
}

// CreatePost ...
func (p *Post) CreatePost(post *app.Post) error {
	if post == nil {
		return errEmpty
	}

	tx := p.DB.MustBegin()

	post.CreatedAt = time.Now().Unix()

	res, err := tx.NamedExec("INSERT INTO posts (creator_id, post_title, post_body, created_at, deleted_at, updated_at) VALUES(:creator_id, :post_title, :post_body, :created_at, :deleted_at, :updated_at)", &post)

	if err != nil && res == nil {
		tx.Rollback()
		return errNotInserted
	}
	tx.Commit()
	return nil
}

// Post ...
func (p *Post) Post(id int64) (*app.Post, error) {
	return nil, nil
}

// Posts ...
func (p *Post) Posts() ([]*app.Post, error) {

	return nil, nil
}

// UpdatePost ...
func (p *Post) UpdatePost(post *app.Post) error {
	return nil
}

// DeletePost ...
func (p *Post) DeletePost(id int64) error {
	return nil
}
