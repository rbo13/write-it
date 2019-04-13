package sql

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/jwtservice"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUserNotInserted      = errors.New("Failed to insert the user")
	errUserUpdate           = errors.New("Failed to updated the user")
	errUserDelete           = errors.New("Failed to delete the user")
	errEmailAlreadyTaken    = errors.New("Email Address is already taken")
	errUsernameAlreadyTaken = errors.New("Username is already taken")
	errNoResultSet          = errors.New("sql: no rows in result set")
	errEmailRequired        = errors.New("Email is required")
	errUsernameRequired     = errors.New("Username is required")
	errMissingCredentials   = errors.New("Email or Password is missing")
	errCredentialsIncorrect = errors.New("Email or Password is invalid")
)

// UserService implements the app.UserService
type UserService interface {
	app.UserService
}

// User implements the UserService interface
type User struct {
	DB        *sqlx.DB
	UserSrvc  *app.User
	TokenAuth *jwtauth.JWTAuth
}

// NewUserSQLService returns the interface that implements the app.UserService
func NewUserSQLService(db *sqlx.DB, jwtService *jwtservice.JWT) UserService {
	return &User{
		DB:        db,
		UserSrvc:  new(app.User),
		TokenAuth: jwtService.TokenAuth,
	}
}

// CreateUser ...
func (u *User) CreateUser(user *app.User) error {
	userRes, err := u.UserByEmail(user.EmailAddress)

	if err != nil && err.Error() != errNoResultSet.Error() {
		return errUserNotInserted
	}

	if userRes != nil {
		return errEmailAlreadyTaken
	}

	userRes, err = u.UserByUsername(user.Username)

	if err != nil && err.Error() != errNoResultSet.Error() {
		return errUserNotInserted
	}

	if userRes != nil {
		return errUsernameAlreadyTaken
	}

	tx := u.DB.MustBegin()

	if userRes == nil {
		user.CreatedAt = time.Now().Unix()
		user.Password = hashPassword(user.Password)

		if user.UserType == "" {
			user.UserType = "reader"
		}

		res, err := tx.NamedExec("INSERT INTO users (username, email, password, user_type, created_at, deleted_at, updated_at) VALUES(:username, :email, :password, :user_type, :created_at, :deleted_at, :updated_at)", &user)

		if err != nil && res == nil {
			tx.Rollback()
			return errUserNotInserted
		}
		tx.Commit()
	}

	return nil
}

// User ...
func (u *User) User(id int64) (*app.User, error) {
	user := new(app.User)

	err := u.DB.Get(user, "SELECT * FROM users WHERE id = ? LIMIT 1;", id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserByEmail ...
func (u *User) UserByEmail(email string) (*app.User, error) {

	if email == "" {
		return nil, errEmailRequired
	}

	user := app.User{}

	err := u.DB.Get(&user, "SELECT * FROM users WHERE email = ? LIMIT 1;", email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UserByUsername ...
func (u *User) UserByUsername(username string) (*app.User, error) {

	if username == "" {
		return nil, errUsernameRequired
	}

	user := app.User{}

	err := u.DB.Get(&user, "SELECT * FROM users WHERE username = ? LIMIT 1;", username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login ...
func (u *User) Login(email, password string) (*app.User, error) {
	if email == "" || password == "" {
		return nil, errMissingCredentials
	}

	user := app.User{}

	// We get a user using the email
	err := u.DB.Get(&user, "SELECT password FROM users WHERE email = ? LIMIT 1;", email)

	if err != nil {
		return nil, err
	}

	passwordsEqual := comparePasswords(user.Password, []byte(password))

	if passwordsEqual {
		err = u.DB.Get(&user, "SELECT * FROM users WHERE email = ? AND password = ? LIMIT 1;", email, user.Password)

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// GenerateAuthToken ...
func (u *User) GenerateAuthToken(user *app.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.ID,
		"email":         user.EmailAddress,
		"authenticated": true,
		"created_at":    user.CreatedAt,
	}

	jwtauth.SetExpiryIn(claims, 1*time.Hour)
	jwtauth.SetIssuedNow(claims)

	_, authToken, err := u.TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

// Users ...
func (u *User) Users() ([]*app.User, error) {
	users := []*app.User{}

	err := u.DB.Select(&users, "SELECT * FROM users ORDER BY id DESC;")

	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserPosts returns a slice to pointer of UserPosts.
func (u *User) GetUserPosts(userID int64) ([]*app.UserPosts, error) {
	var userPosts []*app.UserPosts

	query := "SELECT po.`post_title`, po.`post_body`, po.`created_at`, po.`updated_at`, u.`user_type`, u.`email`, u.`username` FROM posts as po, users as u WHERE po.`creator_id` = u.`id` AND u.`id` = ?;"
	err := u.DB.Select(&userPosts, query, userID)

	if err != nil {
		return nil, err
	}

	return userPosts, nil
}

// UpdateUser ...
func (u *User) UpdateUser(user *app.User) error {
	user.UpdatedAt = time.Now().Unix()

	query := fmt.Sprintf("UPDATE users SET username = '%s', email = '%s', password = '%s', user_type = '%s', updated_at = '%d' WHERE id = %d;", user.Username, user.EmailAddress, user.Password, user.UserType, user.UpdatedAt, user.ID)

	log.Println(query)

	tx := u.DB.MustBegin()
	res := tx.MustExec(query)

	if res == nil {
		tx.Rollback()
		return errUserUpdate
	}

	tx.Commit()
	return nil
}

// DeleteUser ...
func (u *User) DeleteUser(id int64) error {
	tx := u.DB.MustBegin()

	res := tx.MustExec("DELETE FROM users WHERE id = ?;", id)

	if res == nil {
		tx.Rollback()
		return errUserDelete
	}

	tx.Commit()
	return nil
}

func hashPassword(rawPassword string) (hashedPassword string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	hashedPassword = string(hash)
	return hashedPassword
}

// ComparePasswords compares the hashed and raw password. Returns boolean if equal.
func comparePasswords(hashedPassword string, rawPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, rawPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
