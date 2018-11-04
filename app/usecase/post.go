package usecase

import (
	"encoding/json"
	"net/http"

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

}

func (p *postUsecase) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (p *postUsecase) Update(w http.ResponseWriter, r *http.Request) {

}

func (p *postUsecase) Delete(w http.ResponseWriter, r *http.Request) {

}
