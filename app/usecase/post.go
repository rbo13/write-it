package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/go-chi/render"
	"github.com/rbo13/write-it/app"
)

type postUsecase struct {
	postService app.PostService
}

type postResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

// NewPost ...
func NewPost(postService app.PostService) app.Handler {
	return &postUsecase{
		postService,
	}
}

func (p *postUsecase) Create(w http.ResponseWriter, r *http.Request) {
	var post app.Post

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		postResp := postResponse{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &postResp)
		return
	}

	post.CreatorID = int64(claims["user_id"].(float64))

	err = json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		postResp := postResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &postResp)
		return
	}

	err = p.postService.CreatePost(&post)

	if err != nil {
		postResp := postResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &postResp)
		return
	}
	postResp := postResponse{
		StatusCode: http.StatusOK,
		Message:    "Post has been created",
		Success:    true,
		Data:       post,
	}

	render.JSON(w, r, &postResp)
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
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err)
	}

	post.ID = postID
	post.UpdatedAt = time.Now().Unix()

	err = json.NewDecoder(r.Body).Decode(&post)

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
