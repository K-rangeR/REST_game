// database contains the code nessesary for interacting with
// games that are stored in the postgres database
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "game_db"
)

var db *sql.DB

func init() {
	var err error
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("Was unable to connect to the database")
		panic(err)
	}
	checkDBConnection()
}

// checkDBConnection attempts to connect the database, panics if
// there was an error
func checkDBConnection() {
	err := db.Ping()
	if err != nil {
		fmt.Println("Was unable to connect to the database")
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
}
