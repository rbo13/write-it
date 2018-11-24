package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/rbo13/write-it/app"
)

type postUsecase struct {
	postService app.PostService
}

// NewPost ...
func NewPost(postService app.PostService) app.Handler {
	return &postUsecase{
		postService,
	}
}

func (p *postUsecase) Create(w http.ResponseWriter, r *http.Request) {
	var post app.Post

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.postService.CreatePost(&post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, post)
}

func (p *postUsecase) Get(w http.ResponseWriter, r *http.Request) {
	posts, err := p.postService.Posts()

	if err != nil {
		render.JSON(w, r, err.Error())
	}
	render.JSON(w, r, posts)
}

func (p *postUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err.Error())
	}

	post, err := p.postService.Post(postID)

	if err != nil {
		render.JSON(w, r, err.Error())
	}

	render.JSON(w, r, post)
}

func (p *postUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var post app.Post
	// postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	// if err != nil {
	// 	render.JSON(w, r, err)
	// }
	postID := chi.URLParam(r, "id")
	log.Println(postID)
	// post.ID = postID
	post.UpdatedAt = time.Now()

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		render.JSON(w, r, err)
	}

	err = p.postService.UpdatePost(&post)

	if err != nil {
		render.JSON(w, r, err)
	}

	render.JSON(w, r, &post)
}

func (p *postUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err)
	}

	err = p.postService.DeletePost(postID)

	if err != nil {
		render.JSON(w, r, err)
	}

	render.JSON(w, r, "Post Successfully Deleted")
}
