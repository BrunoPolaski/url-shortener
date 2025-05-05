package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type mySQLAdapter struct{}

func NewMySQLAdapter() Database {
	return &mySQLAdapter{}
}

func (m *mySQLAdapter) Connect() *sql.DB {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbType := os.Getenv("DB_TYPE")

	fmt.Printf("Connecting to database %s://%s:%s@tcp(%s:%s)/%s\n", dbType, dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	return db
}
