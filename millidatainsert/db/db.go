package db

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5000
	user     = "elvisgasana"
	password = ""
	dbname   = "postgres"
)

var DB *sql.DB

func InitDb() (*sql.DB, error) {
	var err error

	psqlInfo := "host=localhost dbname=postgres sslmode=disable user=elvisgasana password="
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer DB.Close()

	if err != nil {
		return nil, errors.New("failed to load db")
	}
	DB.SetMaxOpenConns(150)
	DB.SetConnMaxIdleTime(4)
	createTables()

	return DB, nil

}

func createTables() {
	createmydb := `
	CREATE TABLE IF NOT EXISTS products(
		productId TEXT NOT NULL,
		country TEXT NOT NULL,
		location TEXT NOT NULL,
		price INTEGER NOT NULL,
		product TEXT NOT NULL,
		registeredAt TEXT NOT NULL,
		shop TEXT NOT NULL,
		type TEXT NOT NULL
	)
	`
	_, err := DB.Exec(createmydb)
	if err != nil {
		panic("Could not create product table")
	}
}
