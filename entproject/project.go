package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	notes = append(notes, Note{NoteID: 3, UserID: 2, Title: "my note 3", Contents: "hi cat note", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 4, UserID: 1, Title: "my note 4", Contents: "hello world", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 5, UserID: 2, Title: "my note 5", Contents: "hi dog", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 6, UserID: 2, Title: "my note 6", Contents: "pup hi note", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 7, UserID: 1, Title: "my note 7", Contents: "hello doggo", DateCreated: time.Now(), DateUpdated: time.Now()})
	notes = append(notes, Note{NoteID: 8, UserID: 2, Title: "my note 8", Contents: "note is world", DateCreated: time.Now(), DateUpdated: time.Now()})

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
	r.HandleFunc("/Notes/Search", anotherSearch).Methods("POST")

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
	//w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT * FROM Note`)
	if err != nil {
		log.Fatal(err)
	}

	//for each row print ln - need to change to html list at some point
	for rows.Next() {
		var (
			noteID      int
			userID      int
			title       string
			contents    string
			dateCreated time.Time
			dateUpdated time.Time
		)
		err = rows.Scan(&noteID, &userID, &title, &contents, &dateCreated, &dateUpdated)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(noteID, userID, title, contents, dateCreated, dateUpdated)
	}

	//Error check
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	//json.NewEncoder(w).Encode(notes)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(users)
	rows, err := db.Query(`SELECT userID, givenName, familyName FROM "User"`)
	if err != nil {
		log.Fatal(err)
	}

	//for each row print ln - need to change to html list
	for rows.Next() {
		var (
			userID     int
			givenName  string
			familyName string
		)
		err = rows.Scan(&userID, &givenName, &familyName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(userID, givenName, familyName)
	}
	//Error check
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
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

	//Prepare query
	query := `INSERT INTO Note (UserID, Title, Contents, DateCreated, DateUpdated) VALUES ($1, $2, $3, $4, $5)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	//Get todays date
	date := time.Now()
	_, err = stmt.Exec(newNote.UserID, newNote.Title, newNote.Contents, date, date)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created new note")
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
	query := `INSERT INTO "User" (GivenName, FamilyName, Password) VALUES ($1, $2, $3) RETURNING UserID;`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	userID := 0
	err = stmt.QueryRow(newUser.GivenName, newUser.FamilyName, newUser.Password).Scan(&userID)
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
		key := len(arr[i].Contents)
		ts := arr[i]
		j := i - 1
		for j >= 0 && key < len(arr[j].Contents) {
			arr[j+1] = arr[j]
			j -= 1
		}
		arr[j+1] = ts
	}
	fmt.Println(arr)
	return arr
}

// var finalvalue Note

// func searchbin(w http.ResponseWriter, r *http.Request, input string, sortednotes []Note) int { //T is the lastname you are searching for
// 	//sortednotes := insertionSort(notes)
// 	low := 0
// 	high := len(sortednotes) - 1
// 	mid := 0
// 	var mid_value Note
// 	//var input Note
// 	//_ = json.NewDecoder(r.Body).Decode(&input)

// 	for low <= high {
// 		mid = low + (high-low)/2     //middle of the list
// 		mid_value = sortednotes[mid] //get item to check if matches with T

// 		if mid_value.Contents == input || strings.Contains(mid_value.Contents, input) {
// 			//json.NewEncoder(w).Encode(mid_value)
// 			addFoundNote(mid_value)
// 			return mid //we have matched the target T

// 		} else if mid_value.Contents < input {
// 			low = mid + 1 //left/lower side of the middle

// 		} else {
// 			high = mid - 1 //right/upper side of the middle
// 		}

// 	}

// 	return -1 //not found
// }

var foundnotes []Note

func addFoundNote(note Note) {

	if len(foundnotes) == 0 {
		foundnotes = append(foundnotes, note)
	} else {
		for index := 0; index < len(foundnotes); index++ {
			if foundnotes[index].Title == note.Title {

				return
			}
		}
		foundnotes = append(foundnotes, note)

	}

}

// func searchPartial(w http.ResponseWriter, r *http.Request) { //T is the lastname you are searching for
// 	foundnotes = nil
// 	sortednotes := insertionSort(notes)
// 	low := 0
// 	high := len(sortednotes) - 1
// 	mid := 0
// 	var mid_value Note
// 	var input Note
// 	_ = json.NewDecoder(r.Body).Decode(&input)

// 	for low <= high {
// 		mid = low + (high-low)/2     //middle of the list
// 		mid_value = sortednotes[mid] //get item to check if matches with T

// 		if mid_value.Contents == input.Contents || (mysearch(mid_value.Contents, input.Contents) == 0) {
// 			addFoundNote(mid_value)
// 			json.NewEncoder(w).Encode(foundnotes)
// 			return
// 			//json.NewEncoder(w).Encode(foundnotes)
// 			//we have matched the target T

// 		} else if (mid_value.Contents < input.Contents) || (mysearch(mid_value.Contents, input.Contents) == -1) {
// 			low = mid + 1 //left/lower side of the middle

// 		} else {
// 			high = mid - 1 //right/upper side of the middle
// 		}
// 	}
// 	json.NewEncoder(w).Encode(foundnotes)
// 	return //not found
// }

//close to working but still skips over some elements
func partialSearch(w http.ResponseWriter, r *http.Request) {
	foundnotes = nil
	sortednotes := insertionSort(notes)
	lowerlow := 0
	higherhigh := len(sortednotes) - 1
	mid := lowerlow + ((higherhigh - lowerlow) >> 1)

	var input Note
	_ = json.NewDecoder(r.Body).Decode(&input)

	foundAllLower := false
	for foundAllLower == false {
		if searchLower(sortednotes, input.Contents, lowerlow, mid) == false {
			foundAllLower = true
		}
	}
	foundAllHigher := false
	for foundAllHigher == false {
		if searchLower(sortednotes, input.Contents, mid+1, higherhigh) == false {
			foundAllHigher = true
		}
	}

	json.NewEncoder(w).Encode(foundnotes)
}

func searchLower(sortednotes []Note, input string, low int, high int) bool {

	for low <= high {
		mid := low + ((high - low) >> 1) //middle of the list
		mid_value := sortednotes[mid]

		if mid_value.Contents == input || (mysearch(mid_value.Contents, input) == 0) {
			addFoundNote(mid_value)

		} // else if (mid_value.Contents < input) || (mysearch(mid_value.Contents, input) == -1) {
		// 	low = mid + 1 //left/lower side of the middle

		// } else {
		// 	high = mid - 1 //right/upper side of the middle
		// }

		if len(sortednotes[mid].Contents) >= len(input) {
			return searchLower(sortednotes, input, low, mid-1)
		} else {
			return searchLower(sortednotes, input, mid+1, high)
		}
	}
	return false //not found
}

func searchHigher(sortednotes []Note, input string, low int, high int) bool {

	for low <= high {
		mid := low + ((high - low) >> 1) //middle of the list
		mid_value := sortednotes[mid]

		if mid_value.Contents == input || (mysearch(mid_value.Contents, input) == 0) {
			addFoundNote(mid_value)

			//return true
		} //else if (mid_value.Contents < input) || (mysearch(mid_value.Contents, input) == -1) {
		// 	low = mid + 1 //left/lower side of the middle

		// } else {
		// 	high = mid - 1 //right/upper side of the middle
		// }

		if len(sortednotes[mid].Contents) > len(input) {
			return searchHigher(sortednotes, input, low, mid-1)
		} else {
			return searchHigher(sortednotes, input, mid+1, high)
		}
	}
	return false //not found
}

func mysearch(txt string, pattern string) int {
	flag := 0
	if !strings.Contains(txt, pattern) {
		flag = -1
	}
	return flag
}

//fully working but not using binary
func anotherSearch(w http.ResponseWriter, r *http.Request) {

	var input Note
	_ = json.NewDecoder(r.Body).Decode(&input)
	foundnotes = nil
	sortednotes := insertionSort(notes)
	for i := 0; i < len(sortednotes); i++ {
		if sortednotes[i].Contents == input.Contents || (mysearch(sortednotes[i].Contents, input.Contents) == 0) {
			addFoundNote(sortednotes[i])
		}
	}
	json.NewEncoder(w).Encode(foundnotes)

}
