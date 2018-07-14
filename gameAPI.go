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
	games = append(games, newGame)
}

// handleGet will search the DB for the specified game title
// if found it will return json containing the game data
func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	for _, game := range games {
		if gameTitle == game.Title {
			json.NewEncoder(w).Encode(&game)
			return
		}
	}
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(&Game{})
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
	for ndx, game := range games {
		if gameTitle == game.Title {
			games = append(games[:ndx], games[ndx+1:]...)
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(404)
}

// handleGetDeveloper will get a list of all the games with
// the specified developer
func handleGetDeveloper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	developer := vars["developer"]
	list := make([]Game, 0)
	for _, game := range games {
		if game.Developer == developer {
			list = append(list, game)
		}
	}
	if err := json.NewEncoder(w).Encode(&list); err != nil {
		w.WriteHeader(500)
	}
}

// handleGetRating will get a list of all the games with the specified rating
func handleGetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rating := vars["rating"]
	list := make([]Game, 0)
	for _, game := range games {
		if game.Rating == rating {
			list = append(list, game)
		}
	}
	if err := json.NewEncoder(w).Encode(&list); err != nil {
		w.WriteHeader(500)
	}
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
