package sql

// Schemas is a function that returns a slice of string
// that contains the create sql syntax
func Schemas() []string {

	return []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id bigint NOT NULL AUTO_INCREMENT,
			username varchar(16),
			email varchar(151),
			password varchar(255),
			created_at bigint,
			updated_at bigint,
			deleted_at bigint,
			PRIMARY KEY (id)
		);`,

		`
		CREATE TABLE IF NOT EXISTS posts (
			id bigint NOT NULL AUTO_INCREMENT,
			creator_id bigint,
			post_title text,
			post_body text,
			created_at bigint,
			updated_at bigint,
			deleted_at bigint,
			PRIMARY KEY (id),
			FOREIGN KEY (creator_id) REFERENCES users(id)
		);`,
	}
}
