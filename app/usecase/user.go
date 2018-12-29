package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/rbo13/write-it/app"
)

type userUsecase struct {
	userService app.UserService
}

type userResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
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
		createResp := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}

		render.JSON(w, r, &createResp)
		return
	}

	err = u.userService.CreateUser(&user)

	if err != nil {
		createResp := userResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}

		render.JSON(w, r, &createResp)
		return
	}

	createResp := userResponse{
		StatusCode: http.StatusOK,
		Message:    "User successfully registered",
		Success:    true,
		Data:       user,
	}

	render.JSON(w, r, &createResp)
}

func (u *userUsecase) Login(w http.ResponseWriter, r *http.Request) {
	var user app.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		loginResp := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}

		render.JSON(w, r, &loginResp)
		return
	}

	userResp, err := u.userService.UserByEmail(user.EmailAddress)

	if err != nil {
		loginResp := userResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}

		render.JSON(w, r, &loginResp)
		return
	}

	loginResp := userResponse{
		StatusCode: http.StatusOK,
		Message:    "Logged in successfully",
		Success:    true,
		Data:       userResp,
	}

	render.JSON(w, r, &loginResp)
}

func (u *userUsecase) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.Users()

	if err != nil {
		render.JSON(w, r, err.Error())
	}
	render.JSON(w, r, users)
}

func (u *userUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err.Error())
	}

	user, err := u.userService.User(userID)

	if err != nil {
		render.JSON(w, r, err.Error())
	}

	render.JSON(w, r, user)
}

func (u *userUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var user app.User
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err)
	}

	user.ID = userID
	user.UpdatedAt = time.Now().Unix()

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		render.JSON(w, r, err)
	}

	err = u.userService.UpdateUser(&user)

	if err != nil {
		render.JSON(w, r, err)
	}

	render.JSON(w, r, &user)
}

func (u *userUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		render.JSON(w, r, err)
	}

	err = u.userService.DeleteUser(userID)

	if err != nil {
		render.JSON(w, r, err)
	}

	render.JSON(w, r, "User Successfully Deleted")
}
