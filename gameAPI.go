package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Game represents some data that is part of a vidio game
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
	games[newGame.Title] = newGame
}

// handleGet will search the DB for the specified game title
// if found it will return json containing the game data
func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	game, ok := games[gameTitle]
	if ok {
		jsonData, err := json.MarshalIndent(&game, "", "\t")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	} else {
		str := fmt.Sprintf("Could not find %s anywhere in the database\n", gameTitle)
		w.Write([]byte(str))
		w.WriteHeader(404)
	}
}

// handleDelete will remove the specified game from the database
func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameTitle := vars["title"]
	if _, ok := games[gameTitle]; ok {
		delete(games, gameTitle)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

// Represents the database for the REST api for now
var games map[string]Game

func main() {
	games = make(map[string]Game)

	// mock data
	games["skyrim"] = Game{"skyrim", "bethesda", "M"}
	games["2k"] = Game{"2k", "2k games", "E"}
	games["FIFA18"] = Game{"FIFA18", "EA", "E"}
	games["fallout3"] = Game{"fallout3", "bethesda", "M"}

	r := mux.NewRouter()
	r.HandleFunc("/gameAPI/add", handleAdd).Methods("POST")
	r.HandleFunc("/gameAPI/get/{title}", handleGet).Methods("GET")
	r.HandleFunc("/gameAPI/delete/{title}", handleDelete).Methods("DELETE")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
