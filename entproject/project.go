package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
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
	r.HandleFunc("/Users/Notes/{UserID}", getUserNotes).Methods("GET")
	r.HandleFunc("/Notes/Create/", createNote) //.Methods("POST")
	r.HandleFunc("/Notes/{NoteID}", updateNote).Methods("PUT")
	r.HandleFunc("/Notes/{NoteID}", deleteNote).Methods("DELETE")
	r.HandleFunc("/Users/Create", createUser) //.Methods("POST")
	r.HandleFunc("/Users", getUsers).Methods("GET")
	r.HandleFunc("/Users/LogIn", logIn)    //.Methods("POST")
	r.HandleFunc("/Notes/Search/", search) //.Methods("POST")
	r.HandleFunc("/Notes/Analyse", analyseNote).Methods("POST")

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
	rows, err := db.Query(`SELECT * FROM Note`)
	if err != nil {
		log.Fatal(err)
	}
	var notes []Note
	var note Note

	//for each row print ln - need to change to html list at some point
	for rows.Next() {

		err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(noteID, userID, title, contents, dateCreated, dateUpdated)
		notes = append(notes, note)
	}

	//Error check
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(notes)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	t, err := template.ParseFiles("entproject\\UserList.html")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT userID, givenName, familyName FROM "User"`)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	var user User
	//for each row print ln - need to change to html list
	for rows.Next() {

		err = rows.Scan(&user.UserID, &user.GivenName, &user.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(userID, givenName, familyName)

		users = append(users, user)
	}
	//Error check
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// json.NewEncoder(w).Encode(users)

	err = t.Execute(w, users)
	if err != nil {
		log.Fatal(err)
	}

}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rows, err := db.Query(`SELECT * FROM note WHERE note.noteid = ` + params["NoteID"])
	if err != nil {
		log.Fatal(err)
	}
	var note Note
	for rows.Next() {

		err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}

	}
	json.NewEncoder(w).Encode(note)

}

func getUserNotes(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	t, err := template.ParseFiles("entproject\\userhome.html")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT * FROM note WHERE note.userid = ` + params["UserID"])
	if err != nil {
		log.Fatal(err)
	}

	var userNotes []Note
	var note Note

	for rows.Next() {

		err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}
		userNotes = append(userNotes, note)
	}

	err = t.Execute(w, userNotes)
	if err != nil {
		log.Fatal(err)
	}
	//json.NewEncoder(w).Encode(userNotes)
}

//Create a note
func createNote(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	//_ = json.NewDecoder(r.Body).Decode(&newNote)
	t, err := template.ParseFiles("entproject\\createnote.html")

	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		var newNote Note

		newNote.UserID = 1 //need to get from cookie
		newNote.Title = r.FormValue("title")
		newNote.Contents = r.FormValue("content")
		date := time.Now()
		newNote.DateCreated = date
		newNote.DateUpdated = date

		//Prepare query
		query := `INSERT INTO Note (UserID, Title, Contents, DateCreated, DateUpdated) VALUES ($1, $2, $3, $4, $5)`
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newNote.UserID, newNote.Title, newNote.Contents, newNote.DateCreated, newNote.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Created new note")
	//json.NewEncoder(w).Encode(newNote)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newNote Note
	_ = json.NewDecoder(r.Body).Decode(&newNote)

	query := `UPDATE Note SET userid= $1, title = $2, contents = $3, dateupdated = $4 WHERE Note.noteid =` + params["NoteID"]
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	//Get todays date
	date := time.Now()
	_, err = stmt.Exec(newNote.UserID, newNote.Title, newNote.Contents, date)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(newNote)

	// for index, item := range notes {
	// 	if strconv.Itoa(item.NoteID) == params["NoteID"] {
	// 		notes = append(notes[:index], notes[index+1:]...)
	// 		var newNote Note
	// 		_ = json.NewDecoder(r.Body).Decode(&newNote)

	// 		newNoteID, err := strconv.Atoi(params["NoteID"])
	// 		if err == nil {
	// 			newNote.NoteID = newNoteID
	// 			notes = append(notes, newNote)
	// 			json.NewEncoder(w).Encode(newNote)
	// 		}
	// 		return
	// 	}
	// }
	// json.NewEncoder(w).Encode(notes)

}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	_, err := db.Exec(`DELETE FROM note WHERE note.noteid = ` + params["NoteID"])

	if err != nil {
		log.Fatal(err)
	}

	// for index, item := range notes {
	// 	if strconv.Itoa(item.NoteID) == params["NoteID"] {
	// 		notes = append(notes[:index], notes[index+1:]...)
	// 		break
	// 	}
	// }
	// json.NewEncoder(w).Encode(notes)
}

// Creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Content-Type", "application/json")
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
	json.NewEncoder(w).Encode(newUser)*/
	t, err := template.ParseFiles("entproject\\createaccount.html")
	if err != nil {
		log.Fatal(err)
	}

	//When submitted
	if r.Method == "POST" {
		var newUser User
		//Assign input values to newUser
		newUser.GivenName = r.FormValue("givenName")
		newUser.FamilyName = r.FormValue("familyName")
		newUser.Password = r.FormValue("password")

		//Prepare query to insert into DB
		query := `INSERT INTO "User" (GivenName, FamilyName, Password) VALUES ($1, $2, $3) RETURNING UserID;`
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}
		//Used to return userID so we can display it to the user
		userID := 0
		err = stmt.QueryRow(newUser.GivenName, newUser.FamilyName, newUser.Password).Scan(&userID)
		if err != nil {
			log.Fatal(err)
		}
		newUser.UserID = userID
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

//Check ID exists in db
func checkUserID(id int) bool {
	var userID int

	query := `SELECT UserID FROM "User" WHERE UserID = $1;`

	//Prepare query
	userIDCheck, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	err = userIDCheck.QueryRow(id).Scan(&userID)

	//if rows are emtpy, no matching userid
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	return true
}

//Check password and userid matches and exist in db
func checkPassword(password string, userID int) bool {
	var newpass string

	query := `SELECT Password FROM "User" WHERE Password = $1 and UserID = $2`

	//Prepare query
	passwordCheck, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	err = passwordCheck.QueryRow(password, userID).Scan(&newpass)

	//if rows are emtpy, no matching password
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	return true
}

func logIn(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	//_ = json.NewDecoder(r.Body).Decode(&details)
	t, err := template.ParseFiles("entproject\\logintemplate.html")

	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		var logUser User
		//convert input to int
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Fatal(err)
		}
		//set input data to details
		logUser.UserID = id
		logUser.Password = r.FormValue("password")
		//log.Println(logUser) //Checking

		if checkUserID(logUser.UserID) {
			if checkPassword(logUser.Password, logUser.UserID) {
				log.Println("Logged in")
				cookie, err := r.Cookie("logged-in")
				if err == http.ErrNoCookie {
					cookie = &http.Cookie{
						Name:  "logged-in",
						Value: strconv.Itoa(logUser.UserID),
					}
				}

				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)

				//direct to user notes?
			} else {
				log.Println("Failed log in") //http error instead?
				return
			}
		} else {
			log.Println("Failed log in") //http error instead?
			return
		}
	} else {
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func checkLoggedIn() {

}

/*func logIn(w http.ResponseWriter, r *http.Request) {
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

}*/

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
		if searchHigher(sortednotes, input.Contents, mid+1, higherhigh) == false {
			foundAllHigher = true
		}
	}

	json.NewEncoder(w).Encode(foundnotes)
}

func searchLower(sortednotes []Note, input string, low int, high int) bool {

	for low <= high {
		mid := low + ((high - low) >> 1) //middle of the list
		mid_value := sortednotes[mid]

		if mid_value.Contents == input || (contains(mid_value.Contents, input) == 0) {
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

		if mid_value.Contents == input || (contains(mid_value.Contents, input) == 0) {
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

func contains(txt string, pattern string) int {
	flag := 0
	if !strings.Contains(txt, pattern) {
		flag = -1
	}
	return flag
}

//fully working but not using binary
func search(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("entproject\\searchedNotes.html")
	if err != nil {
		log.Fatal(err)
	}
	var searchnotes []Note
	if r.Method == "POST" {
		var input string

		input = r.FormValue("search")

		var note Note

		fmt.Println(input)

		rows, err := db.Query("SELECT * FROM Note WHERE note.contents LIKE " + "'%" + input + "%'")
		if err != nil {
			log.Fatal(err)
		}

		//for each row print ln - need to change to html list at some point
		for rows.Next() {

			err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(noteID, userID, title, contents, dateCreated, dateUpdated)
			searchnotes = append(searchnotes, note)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

	}

	err = t.Execute(w, searchnotes)
	if err != nil {
		log.Fatal(err)
	}

}

func analyseNote(w http.ResponseWriter, r *http.Request) {
	count := 0
	// var input Note
	// _ = json.NewDecoder(r.Body).Decode(&input)
	// foundnotes = nil
	// sortednotes := insertionSort(notes)
	// for i := 0; i < len(sortednotes); i++ {
	// 	if sortednotes[i].Contents == input.Contents || (contains(sortednotes[i].Contents, input.Contents) == 0) {

	// 		count++
	// 	}
	// }

	var input Note
	_ = json.NewDecoder(r.Body).Decode(&input)

	rows, err := db.Query("SELECT * FROM Note WHERE note.contents LIKE " + "'%" + input.Contents + "%'")
	if err != nil {
		log.Fatal(err)
	}

	//for each row print ln - need to change to html list at some point
	for rows.Next() {

		// err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//fmt.Println(noteID, userID, title, contents, dateCreated, dateUpdated)
		count++
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(count)

}
