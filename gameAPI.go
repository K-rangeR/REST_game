package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Game represents some data that is part of a video game
type Game struct {
	Title     string `json:"title"`
	Developer string `json:"developer"`
	Rating    string `json:"rating"`
}

// handleAdd will get the json from the request, convert it
// to a Game struct and store it in memory
func handleAdd(w http.ResponseWriter, r *http.Request) {
	bodySize := r.ContentLength
	bodyData := make([]byte, bodySize)
	r.Body.Read(bodyData)
	var newGame Game
	err := json.Unmarshal(bodyData, &newGame)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(500)
	}
	if err = newGame.addGame(); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(500)
	}
}

// handleGet will search the DB for the specified game title
// if found it will return json containing the game data
func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	game, err := getGameByTitle(gameTitle)
	if err != nil {
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&Game{})
		w.WriteHeader(404)
	}
	json.NewEncoder(w).Encode(&game)
}

// handleUpdate will update the data on an existing game
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	for i, game := range games {
		if gameTitle == game.Title {
			var updatedGame Game
			bodySize := r.ContentLength
			bodyData := make([]byte, bodySize)
			r.Body.Read(bodyData)
			err := json.Unmarshal(bodyData, &updatedGame)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			games[i] = updatedGame
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(404)
}

// handleDelete will remove the specified game from the database
func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	err := deleteGame(gameTitle)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}

// handleGetDeveloper will get a list of all the games with
// the specified developer
func handleGetDeveloper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	developer := vars["developer"]
	games1, err := getGamesByDeveloper(developer) // change back to games
	if err != nil {
		w.WriteHeader(404)
		return
	}
	json.NewEncoder(w).Encode(&games1)
}

// handleGetRating will get a list of all the games with the specified rating
func handleGetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rating := vars["rating"]
	games1, err := getGamesWithRating(rating)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	json.NewEncoder(w).Encode(&games1)
}

// Represents the database for the REST api for now
var games []Game

func main() {
	games = make([]Game, 10)

	// mock data
	games[0] = Game{"skyrim", "bethesda", "M"}
	games[1] = Game{"2k", "2k games", "E"}
	games[2] = Game{"FIFA18", "EA", "E"}
	games[3] = Game{"fallout3", "bethesda", "M"}

	r := mux.NewRouter()
	r.HandleFunc("/gameAPI/add", handleAdd).Methods("POST")
	r.HandleFunc("/gameAPI/{title}", handleGet).Methods("GET")
	r.HandleFunc("/gameAPI/{title}", handleUpdate).Methods("PUT")
	r.HandleFunc("/gameAPI/{title}", handleDelete).Methods("DELETE")
	r.HandleFunc("/gameAPI/developer/{developer}", handleGetDeveloper).Methods("GET")
	r.HandleFunc("/gameAPI/rating/{rating}", handleGetRating).Methods("GET")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
