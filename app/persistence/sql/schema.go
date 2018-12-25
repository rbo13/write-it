package sql

// MigrateSchema is a function that creates the tables to the database.
func MigrateSchema() string {
	return `
		CREATE TABLE users (
			ID int NOT NULL AUTO_INCREMENT
			username varchar(16),
			email varchar(151)
			password varchar(255)
			created_at datetime
			updated_at datetime
			deleted_at datetime
			PRIMARY KEY (ID)
		);

		CREATE TABLE posts (
			id int NOT NULL AUTO_INCREMENT
			post_title varchar(12)
			post_body text
			PRIMARY KEY (ID),
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)`

	// CreatorID int64      `json:"creator_id" db:"creator_id"`
	// PostTitle string     `json:"post_title" db:"post_title"`
	// PostBody  string     `json:"post_body" db:"post_body"`
}
