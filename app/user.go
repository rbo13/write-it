package app

// User represents the user of our application
type User struct {
	ID           int64  `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	EmailAddress string `json:"email_address" db:"email_address"`
	Password     string `json:"-" db:"password"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
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
