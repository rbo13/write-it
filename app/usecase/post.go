package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/persistence/cache"
	"github.com/rbo13/write-it/app/response"
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
		config := response.Configure(err.Error(), http.StatusForbidden, nil)
		response.JSONError(w, r, config)
		return
	}

	post.CreatorID = int64(claims["user_id"].(float64))

	err = json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	err = p.postService.CreatePost(&post)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post created successfully", http.StatusOK, post)
	response.JSONOK(w, r, config)
	return
}

func (p *postUsecase) Get(w http.ResponseWriter, r *http.Request) {

	// get from cache first
	var posts []*app.Post
	cacheKey = "getAllPosts"
	mem := BootMemcached()

	data, err := cache.Get(mem, cacheKey)
	if err == nil && data != "" {
		// val, err := postsUnmarshaler(data, posts)
		err = cache.Unmarshal(data, &posts)

		if err != nil {
			config := response.Configure(err.Error(), http.StatusInternalServerError, nil)
			response.JSONError(w, r, config)
		}

		if err == nil {
			config := response.Configure("Post successfully retrieved", http.StatusOK, map[string]interface{}{
				"posts":  posts,
				"cached": true,
			})
			response.JSONOK(w, r, config)
		}

		return
	}

	posts, err = p.postService.Posts()

	if err != nil {
		config := response.Configure(err.Error(), http.StatusInternalServerError, nil)
		response.JSONError(w, r, config)
		return
	}

	// save to cache
	err = StoreToCache(mem, posts, cacheKey)
	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Posts successfully retrieved", http.StatusOK, map[string]interface{}{
		"posts":  posts,
		"cached": false,
	})

	response.JSONOK(w, r, config)
}

func (p *postUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	// get from cache
	var post *app.Post
	cacheKey = chi.URLParam(r, "id")
	mem := BootMemcached()

	data, err := cache.Get(mem, cacheKey)
	if err == nil && data != "" {
		// err = json.Unmarshal([]byte(data), &post)
		// val, err := Unmarshaler(data, post)
		err = cache.Unmarshal(data, &post)

		if err != nil {
			config := response.Configure(err.Error(), http.StatusInternalServerError, nil)
			response.JSONError(w, r, config)
		}

		// assert the type since we return an interface{}.

		// if val != nil {
		// 	post = val.(*app.Post)
		// }

		if post != nil && err == nil {
			config := response.Configure("Post successfully retrieved", http.StatusOK, map[string]interface{}{
				"post":   post,
				"cached": true,
			})
			response.JSONOK(w, r, config)
		}

		return
	}

	post, err = p.postService.Post(postID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	// save to cache
	err = StoreToCache(mem, post, cacheKey)
	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post successfully retrieved", http.StatusOK, map[string]interface{}{
		"post":   post,
		"cached": false,
	})
	response.JSONOK(w, r, config)
}

func (p *postUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var post app.Post
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	check(err, w, r)

	_, claims, err := jwtauth.FromContext(r.Context())
	check(err, w, r)

	// find a user by the given id
	postFetchRes, err := p.postService.Post(postID)
	check(err, w, r)

	userID := int64(claims["user_id"].(float64))
	if postFetchRes.CreatorID != userID {
		config := response.Configure("Cannot update other Post", http.StatusForbidden, nil)
		response.JSONError(w, r, config)
		return
	}

	post.ID = postFetchRes.ID
	post.CreatorID = int64(claims["user_id"].(float64))
	post.CreatedAt = postFetchRes.CreatedAt

	err = json.NewDecoder(r.Body).Decode(&post)

	check(err, w, r)

	err = p.postService.UpdatePost(&post)

	check(err, w, r)

	config := response.Configure("Post Successfully Updated", http.StatusOK, post)
	response.JSONOK(w, r, config)
}

func (p *postUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	err = p.postService.DeletePost(postID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post Successfully Deleted", http.StatusOK, nil)
	response.JSONOK(w, r, config)
}

func check(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}
	return
}

func postsUnmarshaler(data string, val []*app.Post) ([]*app.Post, error) {
	err := json.Unmarshal([]byte(data), &val)

	if err != nil {
		return nil, err
	}

	return val, nil
}

// Unmarshaler handles the unmarshaling of data
func Unmarshaler(data string, val interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(data), &val)
	if err != nil {
		return nil, err
	}

	return val, nil
}
