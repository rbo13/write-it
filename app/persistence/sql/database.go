package sql

import (
	"fmt"

	// Import the mysql database
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// New ...
func New(dbName string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s ssqlmode=%s", "root", dbName, "disable")
	sqlDB, err := sqlx.Connect("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}
