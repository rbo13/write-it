package sql

import (
	"errors"
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
	errNoResultSet          = errors.New("sql: no rows in result set")
	errEmailRequired        = errors.New("Email is required")
	errMissingCredentials   = errors.New("Email or Password is missing")
	errCredentialsIncorrect = errors.New("Email or Password is invalid")
)

// Servicer ...
type Servicer interface {
	app.UserService
	app.PostService
}

// Service ...
type Service struct {
	DB        *sqlx.DB
	UserSrvc  *app.User
	PostSrvc  *app.Post
	TokenAuth *jwtauth.JWTAuth
}

// NewSQLService ...
func NewSQLService(db *sqlx.DB, jwtService *jwtservice.JWT) Servicer {
	return &Service{
		DB:        db,
		UserSrvc:  new(app.User),
		PostSrvc:  new(app.Post),
		TokenAuth: jwtService.TokenAuth,
	}
}

// CreateUser ...
func (s *Service) CreateUser(user *app.User) error {
	userRes, err := s.UserByEmail(user.EmailAddress)

	if err != nil && err.Error() != errNoResultSet.Error() {
		return errUserNotInserted
	}

	if userRes != nil {
		return errEmailAlreadyTaken
	}

	tx := s.DB.MustBegin()

	if userRes == nil {
		user.CreatedAt = time.Now().Unix()
		user.Password = hashPassword(user.Password)

		res, err := tx.NamedExec("INSERT INTO users (username, email, password, created_at, deleted_at, updated_at) VALUES(:username, :email, :password, :created_at, :deleted_at, :updated_at)", &user)

		if err != nil && res == nil {
			tx.Rollback()
			return errUserNotInserted
		}
		tx.Commit()
	}

	return nil
}

// User ...
func (s *Service) User(id int64) (*app.User, error) {
	user := new(app.User)

	err := s.DB.Get(user, "SELECT * FROM users WHERE id = ? LIMIT 1;", id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserByEmail ...
func (s *Service) UserByEmail(email string) (*app.User, error) {

	if email == "" {
		return nil, errEmailRequired
	}

	user := app.User{}

	err := s.DB.Get(&user, "SELECT * FROM users WHERE email = ? LIMIT 1;", email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login ...
func (s *Service) Login(email, password string) (*app.User, error) {
	if email == "" || password == "" {
		return nil, errMissingCredentials
	}

	user := app.User{}

	// We get a user using the email
	err := s.DB.Get(&user, "SELECT password FROM users WHERE email = ? LIMIT 1;", email)

	if err != nil {
		return nil, err
	}

	passwordsEqual := s.ComparePasswords(user.Password, []byte(password))

	if passwordsEqual {
		err = s.DB.Get(&user, "SELECT * FROM users WHERE email = ? AND password = ? LIMIT 1;", email, user.Password)

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// GenerateAuthToken ...
func (s *Service) GenerateAuthToken(user *app.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.ID,
		"email":         user.EmailAddress,
		"authenticated": true,
		"created_at":    user.CreatedAt,
	}

	jwtauth.SetExpiryIn(claims, 1*time.Hour)
	jwtauth.SetIssuedNow(claims)

	_, authToken, err := s.TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

// Users ...
func (s *Service) Users() ([]*app.User, error) {
	users := []*app.User{}

	err := s.DB.Select(&users, "SELECT * FROM users ORDER BY id DESC;")

	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser ...
func (s *Service) UpdateUser(user *app.User) error {
	user.UpdatedAt = time.Now().Unix()

	tx := s.DB.MustBegin()
	res := tx.MustExec("UPDATE Customers SET username = ?, email = ?, password = ?, updated_at = ?, WHERE id = ?;", user.Username, user.EmailAddress, user.Password, user.UpdatedAt, user.ID)

	if res == nil {
		tx.Rollback()
		return errUserUpdate
	}

	tx.Commit()
	return nil
}

// DeleteUser ...
func (s *Service) DeleteUser(id int64) error {
	tx := s.DB.MustBegin()

	res := tx.MustExec("DELETE FROM users WHERE id = $1;", id)

	if res == nil {
		tx.Rollback()
		return errUserDelete
	}

	tx.Commit()
	return nil
}

// CreatePost ...
func (s *Service) CreatePost(post *app.Post) error {
	return nil
}

// Post ...
func (s *Service) Post(id int64) (*app.Post, error) {
	return nil, nil
}

// Posts ...
func (s *Service) Posts() ([]*app.Post, error) {

	return nil, nil
}

// UpdatePost ...
func (s *Service) UpdatePost(post *app.Post) error {
	return nil
}

// DeletePost ...
func (s *Service) DeletePost(id int64) error {
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

// ComparePasswords compares the hashed and raw password.
// Returns boolean if equal
func (Service) ComparePasswords(hashedPassword string, rawPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, rawPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
