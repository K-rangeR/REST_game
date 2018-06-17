package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
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
	gameTitle := path.Base(r.URL.Path)
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
	gameTitle := path.Base(r.URL.Path)
	if _, ok := games[gameTitle]; ok {
		delete(games, gameTitle)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

var games map[string]Game

func main() {
	games = make(map[string]Game)
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/gameAPI/add", handleAdd)
	http.HandleFunc("/gameAPI/get/", handleGet)
	http.HandleFunc("/gameAPI/delete/", handleDelete)
	server.ListenAndServe()
}
