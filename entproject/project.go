package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Note struct {
	NoteID      int       `json: noteID`
	UserID      int       `json: userID`
	Title       string    `json: title`
	Contents    string    `json: contents`
	DateCreated time.Time `json: dateCreated`
	DateUpdated time.Time `json dateUpdated`
}

type User struct {
	UserID     int    `json: userID`
	GivenName  string `json: givenName`
	FamilyName string `json: familyName`
	Password   string `json: password`
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
var access []NoteAccess
var db *sql.DB

func main() {
	//Router
	r := mux.NewRouter()

	//mock data
	//mock users
	users = append(users, User{UserID: 1, GivenName: "John", FamilyName: "Snow", Password: "hello123"})
	users = append(users, User{UserID: 2, GivenName: "Bob", FamilyName: "Williams", Password: "hi"})
	//mock notes
	notes = append(notes, Note{NoteID: 1, UserID: 1, Title: "my note", Contents: "hi this is a note", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 2, UserID: 1, Title: "my note 2", Contents: "note2", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 3, UserID: 2, Title: "my note 3", Contents: "hi this is a note2", DateCreated: time.Now(), DateUpdated: time.Now()})

	//set up db
	setupDB()
	defer db.Close()
	//Route Handlers
	r.HandleFunc("/Notes", getNotes).Methods("GET")
	r.HandleFunc("/Notes/{NoteID}", getNote).Methods("GET")
	r.HandleFunc("/Notes", createNote).Methods("POST")
	r.HandleFunc("/Notes/{NoteID}", updateNote).Methods("PUT")
	r.HandleFunc("/Notes/{NoteID}", deleteNote).Methods("DELETE")
	r.HandleFunc("/Users/CreateUser", createUser).Methods("POST")
	r.HandleFunc("/Users", getUsers).Methods("GET")
	r.HandleFunc("/Users/LogIn", logIn).Methods("POST")
	r.HandleFunc("/Notes/Search", searchPartial).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func openDB() (db *sql.DB) {
	//Opens database called "EnterpriseNoteApp"
	db, err := sql.Open("postgres", "user=postgres password=password dbname=EnterpriseNoteApp sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func setupDB() {
	//Open db from setupDB file
	db = openDB()

	//Create queries
	createUserTableQuery := `CREATE TABLE IF NOT EXISTS "User"(
		UserID SERIAL PRIMARY KEY,
		GivenName VARCHAR(30),
		FamilyName VARCHAR(30),
		Password VARCHAR(30)
	);`

	createNoteTableQuery := `CREATE TABLE IF NOT EXISTS Note(
		NoteID SERIAL PRIMARY KEY,
		UserID INT,
		Title VARCHAR(30),
		Contents VARCHAR(1000),
		DateCreated DATE,
		DateUpdated DATE,
		FOREIGN KEY (UserID) REFERENCES "User"(UserID)
		);`

	//Execute queries
	_, err := db.Exec(createUserTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createNoteTableQuery)
	if err != nil {
		log.Fatal(err)
	}
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

func getUserNotes(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")
	var userNotes []Note
	for _, item := range notes {
		if item.UserID == user.UserID {

			userNotes = append(userNotes, item)

		}
	}
	json.NewEncoder(w).Encode(userNotes)
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

// Creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User

	_ = json.NewDecoder(r.Body).Decode(&newUser)

	//Return UserID when inserting to let user know their ID
	stmt := `INSERT INTO "User" (GivenName, FamilyName, Password) VALUES ($1, $2, $3) RETURNING UserID;`
	userID := 0
	err := db.QueryRow(stmt, newUser.GivenName, newUser.FamilyName, newUser.Password).Scan(&userID)
	if err != nil {
		log.Fatal(err)
	}
	//Checking
	fmt.Println("Created new user with ID", userID)
	json.NewEncoder(w).Encode(newUser)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var details User
	_ = json.NewDecoder(r.Body).Decode(&details)

	for _, item := range users {

		if item.UserID == details.UserID && item.Password == details.Password {
			getUserNotes(w, r, item)
			return
		}
	}
	fmt.Println("Invalid username or password")

}

// func shareNote(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var noteAccess NoteAccess
// 	_ = json.NewDecoder(r.Body).Decode(&noteAccess)

// 	access.append(access, noteAccess)
// }

// func getSharedNotes(w http.ResponseWriter, r *http.Request) {
// 	var userNotes []Note
// 	for _, item := range access {
// 		if item.UserID == user.UserID {

// 			userNotes = append(userNotes, item)

// 		}
// 	}
// 	json.NewEncoder(w).Encode(userNotes)
// }
func insertionSort(arr []Note) []Note {
	for i := 1; i < len(arr); i++ {
		key := arr[i].Contents
		ts := arr[i]
		j := i - 1
		for j >= 0 && key < arr[j].Contents {
			arr[j+1] = arr[j]
			j -= 1
		}
		arr[j+1] = ts
	}
	return arr
}

var finalvalue Note

func search(input string, sortednotes []Note) int { //T is the lastname you are searching for
	//sortednotes := insertionSort(notes)
	low := 0
	high := len(notes) - 1
	mid := 0
	var mid_value Note
	//var input Note
	//_ = json.NewDecoder(r.Body).Decode(&input)

	for low <= high {
		mid = low + (high-low)/2     //middle of the list
		mid_value = sortednotes[mid] //get item to check if matches with T

		if mid_value.Contents == input {
			//json.NewEncoder(w).Encode(mid_value)
			finalvalue = mid_value
			return mid //we have matched the target T

		} else if mid_value.Contents < input {
			low = mid + 1 //left/lower side of the middle

		} else {
			high = mid - 1 //right/upper side of the middle
		}
	}

	return -1 //not found
}

func searchPartial(w http.ResponseWriter, r *http.Request) { //T is the lastname you are searching for
	sortednotes := insertionSort(notes)
	low := 0
	high := len(notes) - 1
	mid := 0
	var mid_value Note
	var input Note
	_ = json.NewDecoder(r.Body).Decode(&input)

	for low <= high {
		mid = low + (high-low)/2     //middle of the list
		mid_value = sortednotes[mid] //get item to check if matches with T

		if (mid_value.Contents == input.Contents) || (search(mid_value.Contents, sortednotes) == 0) {
			json.NewEncoder(w).Encode(finalvalue)
			return //we have matched the target T

		} else if (mid_value.Contents < input.Contents) || (search(mid_value.Contents, sortednotes) == -1) {
			low = mid + 1 //left/lower side of the middle

		} else {
			high = mid - 1 //right/upper side of the middle
		}
	}

	return //not found
}
