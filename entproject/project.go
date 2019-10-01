package main

import (
	"encoding/json"
	"log"
	"net/http"

	"strconv"
	"time"

	"github.com/gorilla/mux"
	//_ "github.com/lib/pq"
)

type Note struct {
	NoteID      int       `json: noteID`
	Title       string    `json: title`
	Contents    string    `json: contents`
	DateCreated time.Time `json: dateCreated`
	DateUpdated time.Time `json dateUpdated`
}

type User struct {
	UserID     int    `json: userID`
	GivenName  string `json: givenName`
	FamilyName string `json: familyName`
}

type NoteAccess struct {
	NoteAccessID int  `json: noteAccessID`
	NoteID       int  `json: noteID`
	UserID       int  `json: userID`
	Read         bool `json: read`
	Write        bool `json: write`
}

var notes []Note

func main() {
	//Router
	r := mux.NewRouter()

	//mock data
	notes = append(notes, Note{NoteID: 1, Title: "my note", Contents: "hi this is a note", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 2, Title: "my note 2", Contents: "hi this is a note2", DateCreated: time.Now(), DateUpdated: time.Now()})

	//Route Handlers
	r.HandleFunc("/Notes", getNotes).Methods("GET")
	r.HandleFunc("/Notes/{NoteID}", getNote).Methods("GET")
	r.HandleFunc("/Notes", createNote).Methods("POST")
	r.HandleFunc("/Notes/{NoteID}", updateNote).Methods("PUT")
	r.HandleFunc("/Notes/{NoteID}", deleteNote).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range notes {
		if strconv.Itoa(item.NoteID) == params["NoteID"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Note{})
}

func createNote(w http.ResponseWriter, r *http.Request) {

}

func updateNote(w http.ResponseWriter, r *http.Request) {

}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range notes {
		if strconv.Itoa(item.NoteID) == params["NoteID"] {
			notes = append(notes[:index], notes[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(notes)
}

//hi
