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

// addGame addes all the game to the database
func (g *Game) addGame() error {
	statement := `insert into games (title, developer, rating) values ($1, $2, $3)`
	_, err := db.Query(statement, g.Title, g.Developer, g.Rating)
	return err
}

// Searchs the database for a game whos title matches the title given
func getGameByTitle(title string) (*Game, error) {
	game := &Game{}
	statement := "select title, developer, rating from games where title = $1"
	err := db.QueryRow(statement, title).Scan(&game.Title, &game.Developer, &game.Rating)
	return game, err
}

// getGamesByDeveloper will get all games from the database that were made
// by the given developer
func getGamesByDeveloper(developer string) ([]Game, error) {
	statement := `select * from games where developer = $1`
	rows, err := db.Query(statement, developer)
	if err != nil {
		return nil, err
	}

	games := make([]Game, 0)
	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.Title, &game.Developer, &game.Rating)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

// getGamesWithRating will get all games with the given rating
func getGamesWithRating(rating string) ([]Game, error) {
	statement := `select * from games where rating = $1`
	rows, err := db.Query(statement, rating)
	if err != nil {
		return nil, err
	}

	games := make([]Game, 0)
	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.Title, &game.Developer, &game.Rating)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

// deleteGame removes the game from the database with the specific title
func deleteGame(title string) error {
	statement := `delete from games where title = $1`
	_, err := db.Exec(statement, title)
	return err
}
