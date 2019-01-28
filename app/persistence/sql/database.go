package sql

import (
	"log"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// const dbName = "writeit"

// DB uses the sqlx library to interact with our sql database.
type DB struct {
	Sqlx *sqlx.DB
}

// New returns the pointer to sqlx.DB struct.
func New(dataSourceName string) (*DB, error) {
	sqlDB, dbErr := sqlx.Connect("mysql", dataSourceName)

	if dbErr != nil {
		return nil, dbErr
	}

	return &DB{
		Sqlx: sqlDB,
	}, nil
}

// Create creates the database if not exists.
func (db *DB) Create(dbName string) {
	db.Sqlx.MustExec("CREATE DATABASE IF NOT EXISTS `" + dbName + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;")
}

// Use selects the given database to operate with.
func (db *DB) Use(dbName string) {
	db.Sqlx.MustExec("USE " + dbName)
}

// Migrate migrates a table.
func (db *DB) Migrate() {
	for _, schema := range Schemas() {
		db.Sqlx.MustExec(schema)
	}
	log.Println("DB Migrated Successfully")
}
