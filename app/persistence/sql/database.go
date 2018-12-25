package sql

import (

	// Import the mysql database
	"fmt"
	"log"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var sqlDB *sqlx.DB
var dbErr error

const dbName = "writeit"

// New returns the pointer to sqlx.DB struct
func New(dataSourceName string) (*sqlx.DB, error) {
	sqlDB, dbErr = sqlx.Connect("mysql", dataSourceName)

	if dbErr != nil {
		return nil, dbErr
	}
	defer sqlDB.Close()

	res := sqlDB.MustExec("CREATE DATABASE IF NOT EXISTS `" + dbName + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;")

	if res == nil {
		return nil, dbErr
	}

	log.Printf("CREATE RESULT: %v\n\n", res)

	res = sqlDB.MustExec("USE " + dbName)

	if res == nil {
		return nil, dbErr
	}

	return sqlDB, nil
}

func createDatabase(dbName string) string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)
}
