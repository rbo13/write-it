package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/persistence/cache"
	"github.com/rbo13/write-it/app/persistence/cache/memcached"
	"github.com/rbo13/write-it/app/response"
)

var (
	cacheKey = ""
)

const (
	errCacheMiss = "memcache: cache miss"
)

type userUsecase struct {
	userService app.UserService
}

// UserResponse represents a user response
type UserResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

type loginResponse struct {
	UserResponse UserResponse `json:"user_response"`
	AuthToken    string       `json:"auth_token"`
}

// BootMemcached returns the pointer to memcached.Memcached to spin up the caching layer.
func BootMemcached() *memcached.Memcached {
	return memcached.New("localhost", "11211", "localhost:11211")
}

// NewUser ...
func NewUser(userService app.UserService) app.UserHandler {
	return &userUsecase{
		userService,
	}
}

func (u *userUsecase) Create(w http.ResponseWriter, r *http.Request) {
	var user app.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	err = u.userService.CreateUser(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully registered", http.StatusOK, user)
	response.JSONOK(w, r, config)
	return
}

func (u *userUsecase) Login(w http.ResponseWriter, r *http.Request) {
	var user app.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusUnprocessableEntity, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	userResp, err := u.userService.Login(user.EmailAddress, user.Password)

	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusNotFound, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	authToken, err := u.userService.GenerateAuthToken(userResp)

	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusBadRequest, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	loginResp := map[string]interface{}{
		"user":       userResp,
		"auth_token": authToken,
	}

	config := response.Configure("Logged in sucessfully", http.StatusOK, loginResp)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || userID <= 0 {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	mem := BootMemcached()
	cacheKey := chi.URLParam(r, "id")
	var userPosts []*app.UserPosts

	// TODO:: Get From Cache
	err = cache.Get(mem, cacheKey, &userPosts)
	if err == nil {

		config := response.Configure("User Posts successfully retrieved", http.StatusOK, map[string]interface{}{
			"user_posts": userPosts,
			"cached":     true,
		})
		response.JSONOK(w, r, config)
		return
	}

	userPosts, err = u.userService.GetUserPosts(userID)

	if err != nil || userPosts == nil {
		config := response.Configure(err.Error(), http.StatusNotFound, userPosts)
		response.JSONError(w, r, config)
		return
	}

	if len(userPosts) > 0 {
		ok, err := cache.Set(mem, cacheKey, userPosts)

		if err != nil && !ok {
			config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
			response.JSONError(w, r, config)
			return
		}
	}

	config := response.Configure("User Posts successfully retrieved", http.StatusOK, map[string]interface{}{
		"user_posts": userPosts,
		"cached":     false,
	})
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Get(w http.ResponseWriter, r *http.Request) {
	mem := BootMemcached()
	cacheKey = "getAllUsers"
	var usrs []app.User

	err := cache.Get(mem, cacheKey, &usrs)
	if err == nil {
		config := response.Configure("Users successfully retrieved", http.StatusOK, map[string]interface{}{
			"users":  usrs,
			"cached": true,
		})
		response.JSONOK(w, r, config)
		return
	}

	users, err := u.userService.Users()
	if err != nil || users == nil {
		config := response.Configure(err.Error(), http.StatusNotFound, users)
		response.JSONError(w, r, config)
		return
	}

	if len(users) > 0 {
		ok, err := cache.Set(mem, cacheKey, users)

		if err != nil && !ok {
			config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
			response.JSONError(w, r, config)
			return
		}
	}

	config := response.Configure("Users successfully retrieved", http.StatusOK, map[string]interface{}{
		"users":  users,
		"cached": false,
	})
	response.JSONOK(w, r, config)
}

func (u *userUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || userID <= 0 {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	var user *app.User
	cacheKey = chi.URLParam(r, "id")
	mem := BootMemcached()

	err = cache.Get(mem, cacheKey, &user)
	if err == nil {
		config := response.Configure("User successfully retrieved", http.StatusOK, map[string]interface{}{
			"user":   user,
			"cached": true,
		})
		response.JSONOK(w, r, config)
		return
	}

	user, err = u.userService.User(userID)
	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	ok, err := cache.Set(mem, cacheKey, user)
	if err != nil && !ok {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully retrieved", http.StatusOK, map[string]interface{}{
		"user":   user,
		"cached": false,
	})
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var user app.User
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	check(err, w, r)

	_, claims, err := jwtauth.FromContext(r.Context())
	check(err, w, r)

	authorID := int64(claims["user_id"].(float64))
	if userID != authorID {
		config := response.Configure("Cannot update other User", http.StatusForbidden, nil)
		response.JSONError(w, r, config)
		return
	}

	// Find a user by the given id
	userResp, err := u.userService.User(userID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	// fill the necessary fields
	// that doesnt need to be updated
	user.ID = userResp.ID
	user.Password = userResp.Password
	user.CreatedAt = userResp.CreatedAt

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	err = u.userService.UpdateUser(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully updated", http.StatusOK, user)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, userID)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	check(err, w, r)

	authID := int64(claims["user_id"].(float64))
	if userID != authID {
		config := response.Configure("Cannot Delete other User", http.StatusForbidden, nil)
		response.JSONError(w, r, config)
		return
	}

	err = u.userService.DeleteUser(userID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully deleted", http.StatusNoContent, nil)
	response.JSONOK(w, r, config)
}

func errorResponse(statusCode uint, message string) (errResponse UserResponse) {
	errResponse = UserResponse{
		StatusCode: statusCode,
		Message:    message,
		Success:    false,
		Data:       nil,
	}

	return errResponse
}
