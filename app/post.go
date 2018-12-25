package app

import (
	"fmt"
	"time"
)

// Post represents the post inside the application
type Post struct {
	ID        int64      `json:"id" db:"id"`
	CreatorID int64      `json:"creator_id" db:"creator_id"`
	PostTitle string     `json:"post_title" db:"post_title"`
	PostBody  string     `json:"post_body" db:"post_body"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// PostService defines the basic service of post
type PostService interface {
	CreatePost(*Post) error
	Post(id int64) (*Post, error)
	Posts() ([]*Post, error)
	UpdatePost(*Post) error
	DeletePost(id int64) error
}

// TableName represents the table name of post
func (Post) TableName() string {
	return "posts"
}

func (p *Post) String() string {
	return fmt.Sprintf("Post{ID: %d, CreatorID: %d, PostTitle: %s, PostBody: %s}", p.ID, p.CreatorID, p.PostTitle, p.PostBody)
}
