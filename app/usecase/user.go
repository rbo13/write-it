package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/rbo13/write-it/app"

	jwt "github.com/dgrijalva/jwt-go"
)

type userUsecase struct {
	userService app.UserService
}

// JWTData represents the jwt for our authentication
type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

type userResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	AuthToken  string      `json:"auth_token"`
}

const jwtSecret = "5f7532af1ee4524945250f694b5bd06f44f9127bfc35924c457dfa7f68356798319d2d2c4bdce5aaee390cdc731585285e1e374fc1a88dcdbe3f21320b602aba"

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
			AuthToken:  "",
		}

		render.JSON(w, r, &loginResp)
		return
	}

	userResp, err := u.userService.Login(user.EmailAddress, user.Password)

	if err != nil {
		loginResp := userResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
			AuthToken:  "",
		}

		render.JSON(w, r, &loginResp)
		return
	}

	// create a signing using JWT
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"user_id":    strconv.Itoa(int(userResp.ID)),
			"user_email": userResp.EmailAddress,
			"created_at": strconv.Itoa(int(userResp.CreatedAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	authToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		loginResp := userResponse{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
			AuthToken:  "",
		}

		render.JSON(w, r, &loginResp)
		return
	}

	loginResp := userResponse{
		StatusCode: http.StatusOK,
		Message:    "Logged in successfully",
		Success:    true,
		Data:       userResp,
		AuthToken:  authToken,
	}

	render.JSON(w, r, &loginResp)
}

func (u *userUsecase) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.Users()

	if err != nil {
		getResponse := userResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &getResponse)
		return
	}

	getResponse := userResponse{
		StatusCode: http.StatusOK,
		Message:    "Users successfully retrieved",
		Success:    true,
		Data:       users,
	}
	render.JSON(w, r, &getResponse)
}

func (u *userUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {

		getByIDResponse := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &getByIDResponse)
		return
	}

	user, err := u.userService.User(userID)

	if err != nil {
		getByIDResponse := userResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &getByIDResponse)
		return
	}

	getByIDResponse := userResponse{
		StatusCode: http.StatusOK,
		Message:    "User successfully retrieved",
		Success:    true,
		Data:       user,
	}
	render.JSON(w, r, &getByIDResponse)
}

func (u *userUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var user app.User
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {

		updateResponse := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &updateResponse)
		return
	}

	user.ID = userID
	user.UpdatedAt = time.Now().Unix()

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {

		updateResponse := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &updateResponse)
		return
	}

	err = u.userService.UpdateUser(&user)

	if err != nil {
		updateResponse := userResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &updateResponse)
		return
	}

	updateResponse := userResponse{
		StatusCode: http.StatusOK,
		Message:    "User successfully updated",
		Success:    true,
		Data:       user,
	}
	render.JSON(w, r, &updateResponse)
}

func (u *userUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		deleteResponse := userResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &deleteResponse)
		return
	}

	err = u.userService.DeleteUser(userID)

	if err != nil {
		deleteResponse := userResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &deleteResponse)
		return
	}

	deleteResponse := userResponse{
		StatusCode: http.StatusNoContent,
		Message:    "User successfully deleted",
		Success:    true,
		Data:       nil,
	}
	render.JSON(w, r, &deleteResponse)
}
