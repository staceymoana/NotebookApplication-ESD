package main

import (
	"encoding/json"
	"log"
	"math/rand"
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
var users []User

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
	r.HandleFunc("/CreateUser", createUser).Methods("POST")
	r.HandleFunc("/Users", getUsers).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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
	w.Header().Set("Content-Type", "application/json")
	var newNote Note
	_ = json.NewDecoder(r.Body).Decode(&newNote)

	newNote.NoteID = rand.Intn(100000)
	notes = append(notes, newNote)
	json.NewEncoder(w).Encode(newNote)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range notes {
		if strconv.Itoa(item.NoteID) == params["NoteID"] {
			notes = append(notes[:index], notes[index+1:]...)
			var newNote Note
			_ = json.NewDecoder(r.Body).Decode(&newNote)

			newNoteID, err := strconv.Atoi(params["NoteID"])
			if err == nil {
				newNote.NoteID = newNoteID
				notes = append(notes, newNote)
				json.NewEncoder(w).Encode(newNote)
			}
			return
		}
	}
	json.NewEncoder(w).Encode(notes)

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

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	_ = json.NewDecoder(r.Body).Decode(&newUser)

	newUser.UserID = rand.Intn(100000)
	users = append(users, newUser)
	json.NewEncoder(w).Encode(newUser)
}
