// database contains the code nessesary for interacting with
// games that are stored in the postgres database
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// dbCredentials contains all the credentials used to connect to a particular
// postgres database. This info is stored in a JSON file.
type dbCredentials struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Dbname string `json:"dbname"`
}

// setUpDatabase will attempt to connect to correct database
// using the credentials in the credentialsFile
func setUpDatabase(credentialsFile string) error {
	var err error

	credentials, err := getDatabaseCredentials(credentialsFile)
	if err != nil {
		fmt.Println("Could not get database login info")
		return err
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		credentials.Host, credentials.Port, credentials.User, credentials.User)

	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("Was unable to connect to the database")
		return err
	}
	checkDBConnection()
	return nil
}

// getDatabaseCredentials loads the database credentials from an
// external json file
func getDatabaseCredentials(credentialsFile string) (dbCredentials, error) {
	var credentials dbCredentials
	jsonFile, err := os.Open(credentialsFile)
	if err != nil {
		return credentials, err
	}
	defer jsonFile.Close()

	err = json.NewDecoder(jsonFile).Decode(&credentials)
	return credentials, err
}

// checkDBConnection will ping the database and return an error
// if it was unable to connect
func checkDBConnection() {
	err := db.Ping()
	if err != nil {
		fmt.Println("Was unable to connect to the database")
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
}

// addGame addes all the game data to the database
func (g *Game) addGame() error {
	statement := `INSERT INTO games (title, developer, rating) VALUES ($1, $2, $3);`
	_, err := db.Query(statement, g.Title, g.Developer, g.Rating)
	return err
}

// getGameByTitle searchs the database for a game whos title matches the title given
func getGameByTitle(title string) (*Game, error) {
	game := &Game{}
	statement := `SELECT title, developer, rating FROM games WHERE title=$1;`
	err := db.QueryRow(statement, title).Scan(&game.Title, &game.Developer, &game.Rating)
	return game, err
}

// getGamesByDeveloper will get all games from the database that were made
// by the given developer
func getGamesByDeveloper(developer string) ([]Game, error) {
	statement := `SELECT * FROM games WHERE developer=$1;`
	return getSliceOfGames(statement, developer)
}

// getGamesWithRating will get all games with the given rating
func getGamesWithRating(rating string) ([]Game, error) {
	statement := `SELECT * FROM games WHERE rating=$1;`
	return getSliceOfGames(statement, rating)
}

// getListOfGames will read rows from the dababase
func getSliceOfGames(statement, value string) ([]Game, error) {
	rows, err := db.Query(statement, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// deleteGame removes the game from the database whos title matches the given title
func deleteGame(title string) error {
	statement := `DELETE FROM games WHERE title=$1;`
	_, err := db.Exec(statement, title)
	return err
}

// updateGame updates the info of a game whos title matches the given title
func (g *Game) updateGame(title string) error {
	statement := `UPDATE games SET title=$2, developer=$3, rating=$4 WHERE title=$1;`
	_, err := db.Exec(statement, title, g.Title, g.Developer, g.Rating)
	return err
}
