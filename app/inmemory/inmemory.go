package inmemory

import (
	"errors"
	"sync"
	"time"

	"github.com/rbo13/write-it/app"
)

var errIDRequired = errors.New("ID is required")

type postService struct {
	mu    *sync.RWMutex
	posts map[int64]*app.Post
}

// NewInMemoryPostService ...
func NewInMemoryPostService() app.PostService {
	return &postService{
		mu:    &sync.RWMutex{},
		posts: map[int64]*app.Post{},
	}
}

func (ps *postService) CreatePost(post *app.Post) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ps.posts[post.ID] = &app.Post{
		ID:        post.ID,
		CreatorID: post.CreatorID,
		PostTitle: post.PostTitle,
		PostBody:  post.PostBody,
		CreatedAt: time.Now(),
		DeletedAt: nil,
	}

	return nil
}

func (ps *postService) Post(id int64) (*app.Post, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if id <= 0 {
		return nil, errIDRequired
	}

	for _, post := range ps.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("Post not found")
}

func (ps *postService) Posts() ([]*app.Post, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var posts []*app.Post

	for _, post := range ps.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *postService) UpdatePost(post *app.Post) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if post.ID <= 0 {
		return errIDRequired
	}

	ps.posts[post.ID] = post

	return nil
}

func (ps *postService) DeletePost(id int64) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if id <= 0 {
		return errIDRequired
	}

	ps.posts[id] = nil

	return nil
}
