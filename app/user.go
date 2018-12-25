package app

import (
	"time"
)

// User represents the user of our application
type User struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	EmailAddress string     `json:"email_address"`
	Password     string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// UserService defines the basic service of user
type UserService interface {
	CreateUser(*User) error
	User(id int64) (*User, error)
	Users() ([]*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
}

// TableName represents the table name of user
func (User) TableName() string {
	return "user"
}
