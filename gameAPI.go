package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Game represents some data that is part of a video game
type Game struct {
	Title     string `json:"title"`
	Developer string `json:"developer"`
	Rating    string `json:"rating"`
}

// handleAdd will get the json from the request, convert it
// to a Game struct and store it in the database
func handleAdd(w http.ResponseWriter, r *http.Request) {
	bodySize := r.ContentLength
	bodyData := make([]byte, bodySize)
	r.Body.Read(bodyData)
	defer r.Body.Close()
	var newGame Game
	err := json.Unmarshal(bodyData, &newGame)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setGameDataCase(&newGame)
	if err = newGame.addGame(); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleGet will search the DB for the specified game title
// if found it will return json containing the game data
func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	fmt.Println(gameTitle)
	gameTitle = strings.ToLower(gameTitle)
	game, err := getGameByTitle(gameTitle)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&game)
	w.WriteHeader(http.StatusOK)
}

// handleUpdate will update the data on an existing game
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	game, err := getGameByTitle(gameTitle)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	defer r.Body.Close()
	json.Unmarshal(body, &game)
	setGameDataCase(game)
	err = game.updateGame(gameTitle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// setGameDataCase will set the games title, developer, and
// rating to lower case before its added to the database
func setGameDataCase(g *Game) {
	strings.ToLower(g.Title)
	strings.ToLower(g.Developer)
	strings.ToUpper(g.Rating)
}

// handleDelete will remove the specified game from the database
func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	gameTitle = strings.ToLower(gameTitle)
	err := deleteGame(gameTitle)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleGetDeveloper will get a list of all the games with
// the specified developer
func handleGetDeveloper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	developer := vars["developer"]
	developer = strings.ToLower(developer)
	games, err := getGamesByDeveloper(developer)
	if err != nil {
		fmt.Println("handle get dev:", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&games)
	w.WriteHeader(http.StatusOK)
}

// handleGetRating will get a list of all the games with the specified rating
func handleGetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rating := vars["rating"]
	rating = strings.ToUpper(rating)
	games, err := getGamesWithRating(rating)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&games) // handle this error
	if err != nil {
		fmt.Println("handleGetRating JSON ecoding error:")
	}
	w.WriteHeader(http.StatusOK)
}

// Allow user to pass in the name of the db dbCredentials file
// as a command line argument, may not be able to use init.
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [database_credentials_file.json]\n", os.Args[0])
		os.Exit(1)
	}
	err := setUpDatabase(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/gameAPI/add", handleAdd).Methods("POST")
	r.HandleFunc("/gameAPI/{title}", handleGet).Methods("GET")
	r.HandleFunc("/gameAPI/{title}", handleUpdate).Methods("PUT")
	r.HandleFunc("/gameAPI/{title}", handleDelete).Methods("DELETE")
	r.HandleFunc("/gameAPI/developer/{developer}", handleGetDeveloper).Methods("GET")
	r.HandleFunc("/gameAPI/rating/{rating}", handleGetRating).Methods("GET")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
