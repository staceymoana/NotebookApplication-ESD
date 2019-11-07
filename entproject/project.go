package main

import (
	"database/sql"
	"fmt"
	"strings"

	//"encoding/json"

	"log"
	"net/http"
	"strconv"
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

type SharedSettings struct {
	SharedSettingsID int
	OwnerID          int
	SharedUserID     int
	Read             bool
	Write            bool
	Name             string
}

//var notes []Note
//var users []User

var db *sql.DB

func main() {
	//Router
	r := mux.NewRouter()

	/*//mock data
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
	notes = append(notes, Note{NoteID: 8, UserID: 2, Title: "my note 8", Contents: "note is world", DateCreated: time.Now(), DateUpdated: time.Now()})*/

	//set up db
	setupDB()
	defer db.Close()
	//Route Handlers
	//r.HandleFunc("/Notes", getNotes).Methods("GET")
	//r.HandleFunc("/Notes/{NoteID}", getNote).Methods("GET")
	r.HandleFunc("/Users/Notes/{UserID}", getUserNotes).Methods("GET")
	r.HandleFunc("/Notes/Create/", createNote)         //.Methods("POST")
	r.HandleFunc("/Notes/Update/{NoteID}", updateNote) //.Methods("PUT")
	r.HandleFunc("/Notes/Delete/{NoteID}", deleteNote) //.Methods("DELETE")
	r.HandleFunc("/Users/Create", createUser)          //.Methods("POST")
	r.HandleFunc("/Users", getUsers).Methods("GET")
	r.HandleFunc("/Users/LogIn", logIn)                  //.Methods("POST")
	r.HandleFunc("/Notes/Search/", search)               //.Methods("POST")
	r.HandleFunc("/Notes/Analyse/{NoteID}", analyseNote) //.Methods("POST")
	r.HandleFunc("/Notes/Share/{NoteID}", shareNote)
	r.HandleFunc("/Notes/ViewAccess/{NoteID}", access)
	r.HandleFunc("/Notes/EditAccess/{NoteID}", editAccess)
	r.HandleFunc("/Notes/CreateSharedSetting/{NoteID}", saveSharedSettingOnNote)
	r.HandleFunc("/Users/Logout", logOut)

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

func setupDB() *sql.DB {
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

	createNoteAccessQuery := `CREATE TABLE IF NOT EXISTS NoteAccess (
		NoteAccessID SERIAL PRIMARY KEY,
		NoteID INT,
		UserID INT,
		Read BOOL,
		Write BOOL,
		FOREIGN KEY (NoteID) REFERENCES Note(NoteID),
		FOREIGN KEY (UserID) REFERENCES "User"(UserID)
	);`

	createSharedSettingsQuery := `CREATE TABLE IF NOT EXISTS SharedSettings  (
		SharedSettingsID SERIAL PRIMARY KEY,
		OwnerID INT, 
		SharedUserID INT,
		Read bool,
		Write bool,
		Name VARCHAR(30),
		FOREIGN KEY (OwnerID) REFERENCES "User"(UserID)
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

	_, err = db.Exec(createNoteAccessQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createSharedSettingsQuery)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//Used for Postman
/*func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT * FROM Note`)
	if err != nil {
		log.Fatal(err)
	}
	var notes []Note
	var note Note

	for rows.Next() {

		err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, note)
	}

	//Error check
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(notes)
}*/

//Gets all users

func getUsers(w http.ResponseWriter, r *http.Request) {
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}
	t, err := template.ParseFiles("templates\\UserList.html")
	if err != nil {
		log.Fatal(err)
	}

	users := getUsersSQL()

	err = t.Execute(w, users)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsersSQL() []User {
	rows, err := db.Query(`SELECT userID, givenName, familyName FROM "User"`)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	var user User

	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.GivenName, &user.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}

//Used for Postman
/*func getNote(w http.ResponseWriter, r *http.Request) {
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

}*/

//Gets all user notes
func getUserNotes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	if cookie.Value == params["UserID"] {
		t, err := template.ParseFiles("templates\\userhome.html")
		if err != nil {
			log.Fatal(err)
		}

		userNotes := getUserNotesSQL(params["UserID"])

		err = t.Execute(w, userNotes)
		if err != nil {
			log.Fatal(err)

		}
	} else {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
	}

}

func getUserNotesSQL(params string) []Note {
	rows, err := db.Query(`SELECT DISTINCT note.noteid,note.userid,note.title,note.contents,note.datecreated,note.dateupdated FROM note LEFT JOIN noteaccess ON note.noteid = noteaccess.noteid WHERE note.userid = ` + params + ` OR (noteaccess.userid = ` + params + ` AND noteaccess.read = true)`)
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
	return userNotes
}

//Creates a note
func createNote(w http.ResponseWriter, r *http.Request) {
	cookie := checkLoggedIn(r)

	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates\\createnote.html")
	if err != nil {
		log.Fatal(err)
	}

	settings := createNoteSelectSQL(cookie.Value)

	if r.Method == "POST" {
		createNoteInsertSQL(cookie.Value, r.FormValue("title"), r.FormValue("content"), r.FormValue("settingSelect"))
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}

	err = t.Execute(w, settings)
	if err != nil {
		log.Fatal(err)

	}
}

func createNoteSelectSQL(userID string) []SharedSettings {
	var settings []SharedSettings
	var setting SharedSettings

	rows, err := db.Query(`SELECT DISTINCT name FROM SharedSettings WHERE OwnerID = ` + userID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&setting.Name)
		if err != nil {
			log.Fatal(err)
		}
		settings = append(settings, setting)
	}
	return settings
}

func createNoteInsertSQL(userID string, title string, content string, selectSetting string) bool {
	var newNote Note
	var err error

	newNote.UserID, err = strconv.Atoi(userID)
	if err != nil {
		log.Fatal(err)
	}
	newNote.Title = title
	newNote.Contents = content
	date := time.Now()
	newNote.DateCreated = date
	newNote.DateUpdated = date

	//Prepare query
	query := `INSERT INTO Note (UserID, Title, Contents, DateCreated, DateUpdated) VALUES ($1, $2, $3, $4, $5) RETURNING NoteID;`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return false
	}

	var noteID int
	err = stmt.QueryRow(newNote.UserID, newNote.Title, newNote.Contents, newNote.DateCreated, newNote.DateUpdated).Scan(&noteID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	newNote.NoteID = noteID

	selectedSetting := selectSetting

	var setting SharedSettings

	rows, err := db.Query(`SELECT SharedSettings.SharedUserID, SharedSettings.Read, SharedSettings.Write FROM SharedSettings WHERE OwnerID = ` + userID + `AND SharedSettings.Name = '` + selectedSetting + `'`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	for rows.Next() {
		err = rows.Scan(&setting.SharedUserID, &setting.Read, &setting.Write)
		if err != nil {
			log.Fatal(err)
			return false
		}
		//settings = append(settings, setting)
		query := `INSERT INTO NoteAccess (NoteID, UserID, Read, Write) VALUES ($1, $2, $3, $4)`
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Fatal(err)
			return false
		}
		_, err = stmt.Exec(noteID, setting.SharedUserID, setting.Read, setting.Write)
		if err != nil {
			log.Fatal(err)
			return false
		}

	}
	return true
}

//Updates a note
func updateNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	cookie := checkLoggedIn(r)

	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	writeValue, note := updateNoteSelectSQL(params["NoteID"])
	id := strconv.Itoa(note.NoteID)

	if id != cookie.Value && writeValue == false {
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}

	t, err := template.ParseFiles("templates\\updatenote.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		updateNoteInsertSQL(r.FormValue("title"), r.FormValue("content"), params["NoteID"])
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}
	err = t.Execute(w, note)
	if err != nil {
		log.Fatal(err)
	}
}

func updateNoteSelectSQL(noteID string) (bool, Note) {
	var writeValue bool
	var note Note

	rows, err := db.Query(`SELECT noteaccess.write From Noteaccess WHERE noteaccess.noteid = ` + noteID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&writeValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	noterow, err := db.Query(`SELECT note.userid,note.title,note.contents FROM note WHERE note.noteid = ` + noteID)

	for noterow.Next() {
		err = noterow.Scan(&note.NoteID, &note.Title, &note.Contents)
		if err != nil {
			log.Fatal(err)
		}
	}
	return writeValue, note
}

func updateNoteInsertSQL(title string, contents string, noteID string) bool {
	var newNote Note

	newNote.Title = title
	newNote.Contents = contents

	query := `UPDATE Note SET title = $1, contents = $2, dateupdated = $3 WHERE Note.noteid =` + noteID
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return false
	}
	//Get todays date
	date := time.Now()
	_, err = stmt.Exec(newNote.Title, newNote.Contents, date)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//Checks whether user logged in is the owner
func isOwner(w http.ResponseWriter, r *http.Request) bool {
	params := mux.Vars(r)
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return false
	}

	userValue := isOwnerSQL(params["NoteID"], cookie.Value)

	if strconv.Itoa(userValue) != cookie.Value {
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
		return false
	}
	return true
}

func isOwnerSQL(noteID string, userID string) int {
	var userValue int

	rows, err := db.Query(`SELECT userid FROM note WHERE note.noteid = ` + noteID + ` AND note.userid = ` + userID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		err = rows.Scan(&userValue)
		if err != nil {
			log.Fatal(err)
		}
	}
	return userValue
}

//Deletes a note
func deleteNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	if isOwner(w, r) {
		deleteNoteSQL(params["NoteID"])
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}
}

func deleteNoteSQL(NoteID string) bool {
	_, err := db.Exec(`DELETE FROM NoteAccess WHERE NoteAccess.noteid = ` + NoteID)
	if err != nil {
		log.Fatal(err)
		return false
	}

	_, err = db.Exec(`DELETE FROM note WHERE note.noteid = ` + NoteID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//Creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates\\createaccount.html")
	if err != nil {
		log.Fatal(err)
	}
	var newUser User
	//When submitted
	if r.Method == "POST" {
		if r.FormValue("givenName") == "" || r.FormValue("familyName") == "" || r.FormValue("password") == "" {
			http.Redirect(w, r, "/Users/Create", http.StatusSeeOther)

		} else {

			newUser = createUserSQL(r.FormValue("givenName"), r.FormValue("familyName"), r.FormValue("password"))
			t2, err := template.ParseFiles("templates\\accountcreated.html")
			if err != nil {
				log.Fatal(err)
			}

			err = t2.Execute(w, newUser)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createUserSQL(givenName string, familyName string, password string) User {
	var newUser User
	//Assign input values to newUser
	newUser.GivenName = givenName
	newUser.FamilyName = familyName
	newUser.Password = password

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
	return newUser
}

//Check password and userid matches and exist in db when a user logs in
func checkPassword(password string, userID int) bool {
	var newpass string

	query := `SELECT Password FROM "User" WHERE Password = $1 and UserID = $2`

	//Prepare query
	passwordCheck, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	err = passwordCheck.QueryRow(password, userID).Scan(&newpass)

	//if rows are empty, no matching password
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	return true
}

//Logs a user in
func logIn(w http.ResponseWriter, r *http.Request) {
	cookie := checkLoggedIn(r)

	if cookie != nil {
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}

	t, err := template.ParseFiles("templates\\logintemplate.html")

	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		idvalue := r.FormValue("id")
		passvalue := r.FormValue("password")

		if idvalue == "" || passvalue == "" {
			http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)

		} else {
			var logUser User
			//convert input to int
			id, err := strconv.Atoi(idvalue)
			if err != nil {
				log.Fatal(err)
			}
			//set input data to details
			logUser.UserID = id
			logUser.Password = passvalue

			if checkPassword(logUser.Password, logUser.UserID) {
				cookie, err := r.Cookie("logged-in")
				if err == http.ErrNoCookie {
					cookie = &http.Cookie{
						Name:  "logged-in",
						Value: strconv.Itoa(logUser.UserID),
						Path:  "/",
					}
				}
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
				return
			}
		}
	} else {
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

}

//Checks whether a user is logged in
func checkLoggedIn(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("logged-in")
	if err == http.ErrNoCookie {
		return nil
	}
	return cookie
}

//fully working but not using binary
func search(w http.ResponseWriter, r *http.Request) {

	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates\\searchedNotes.html")
	if err != nil {
		log.Fatal(err)
	}

	var searchNotes []Note

	if r.Method == "POST" {
		searchNotes = searchSQL(r.FormValue("search"), cookie.Value)
	}

	err = t.Execute(w, searchNotes)
	if err != nil {
		log.Fatal(err)
	}
}

func searchSQL(searchInput string, userid string) []Note {
	var searchNotes []Note
	var input = searchInput

	var note Note

	fmt.Println(input)

	rows, err := db.Query("SELECT DISTINCT note.NoteID, note.UserId, note.title, note.contents, note.datecreated, note.dateupdated FROM note LEFT JOIN noteaccess ON note.noteid = noteaccess.noteid WHERE (note.userid = " + userid + " OR (noteaccess.userid = " + userid + " AND noteaccess.read = true)) AND note.contents LIKE " + "'%" + searchInput + "%'" + " OR note.Title LIKE " + "'%" + searchInput + "%'")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		err = rows.Scan(&note.NoteID, &note.UserID, &note.Title, &note.Contents, &note.DateCreated, &note.DateUpdated)
		if err != nil {
			log.Fatal(err)
		}
		searchNotes = append(searchNotes, note)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return searchNotes
}

//Searches a term and displays a count
func analyseNote(w http.ResponseWriter, r *http.Request) {
	count := 0
	params := mux.Vars(r)

	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates\\analyseNote.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		count = analyseNoteSQL(r.FormValue("search"), params["NoteID"])
	}

	err = t.Execute(w, struct {
		NoteID string
		Count  int
	}{params["NoteID"], count})
	if err != nil {
		log.Fatal(err)
	}

}

func analyseNoteSQL(searchInput string, noteID string) int {
	count := 0
	var input string
	var contents string

	input = searchInput

	rows, err := db.Query("SELECT note.contents FROM Note WHERE Note.Noteid = " + noteID)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	for rows.Next() {
		err = rows.Scan(&contents)
		if err != nil {
			log.Fatal(err)
			return 0
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return 0
	}

	count = strings.Count(contents, input)
	return count
}

//Allows a note to be shared to other users
func shareNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	if isOwner(w, r) {
		t, err := template.ParseFiles("templates\\share.html")
		if err != nil {
			log.Fatal(err)
		}

		if r.Method == "POST" {
			if r.FormValue("userid") != "" {
				shareNoteSQL(r.FormValue("userid"), r.FormValue("readaccess"), r.FormValue("writeaccess"), params["NoteID"])
				http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Notes/Share/"+params["NoteID"], http.StatusSeeOther)
			}

		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func shareNoteSQL(userID string, read string, write string, noteID string) bool {
	var newNoteAccess NoteAccess
	var err error

	newNoteAccess.UserID, err = strconv.Atoi(userID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	newNoteAccess.NoteID, err = strconv.Atoi(noteID)
	if err != nil {
		log.Fatal(err)
		return false
	}

	readValue := read

	if readValue == "on" {
		newNoteAccess.Read = true
	} else {
		newNoteAccess.Read = false
	}

	writeValue := write

	if writeValue == "on" {
		newNoteAccess.Write = true
		newNoteAccess.Read = true
	} else {
		newNoteAccess.Write = false
	}

	//Prepare query
	query := `INSERT INTO NoteAccess (UserID, NoteID, Read, Write) VALUES ($1, $2, $3, $4)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return false
	}

	_, err = stmt.Exec(newNoteAccess.UserID, newNoteAccess.NoteID, newNoteAccess.Read, newNoteAccess.Write)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

//Saves new note access settings
func access(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if isOwner(w, r) {
		t, err := template.ParseFiles("templates\\access.html")
		matches := accessSQL(params["NoteID"])

		err = t.Execute(w, matches)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func accessSQL(noteID string) []NoteAccess {
	matching, err := db.Query(`SELECT na.userid, na.noteid, na.Read, na.Write FROM NoteAccess as na INNER JOIN Note on na.noteid = note.noteid WHERE note.noteid =` + noteID + `AND na.read = true`)
	if err != nil {
		log.Fatal(err)
	}

	var matches []NoteAccess
	var note NoteAccess

	for matching.Next() {

		err = matching.Scan(&note.UserID, &note.NoteID, &note.Read, &note.Write)
		if err != nil {
			log.Fatal(err)
		}
		matches = append(matches, note)
	}
	return matches
}

//Allows a user to edit note access settings
func editAccess(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	if isOwner(w, r) {
		t, err := template.ParseFiles("templates\\editaccess.html")
		if err != nil {
			log.Fatal(err)
		}

		if r.Method == "POST" {
			editAccessSQL(r.FormValue("readaccess"), r.FormValue("writeaccess"), params["NoteID"])
			http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)

		}
	}
}

func editAccessSQL(read string, write string, noteID string) bool {
	var newNoteAccess NoteAccess

	readValue := read

	if readValue == "on" {
		newNoteAccess.Read = true
	} else {
		newNoteAccess.Read = false
	}

	writeValue := write

	if writeValue == "on" {
		newNoteAccess.Write = true
		newNoteAccess.Read = true
	} else {
		newNoteAccess.Write = false
	}

	//Prepare query
	query := `UPDATE NoteAccess SET read = $1, write = $2 WHERE noteaccess.noteid =` + noteID
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return false
	}

	_, err = stmt.Exec(newNoteAccess.Read, newNoteAccess.Write)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

//Allows a user to save certain shared settings and set a name for it
func saveSharedSettingOnNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates\\createSharedSetting.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		saveSharedSettingOnNoteSQL(r.FormValue("settingName"), params["NoteID"])
		http.Redirect(w, r, "/Users/Notes/"+cookie.Value, http.StatusSeeOther)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func saveSharedSettingOnNoteSQL(settingName string, noteID string) bool {
	var setting SharedSettings

	setting.Name = settingName

	rows, err := db.Query(`SELECT n.userid as "owner", na.userid, na.read, na.write FROM NoteAccess as na INNER JOIN Note as n ON na.Noteid = n.noteid WHERE N.noteid = ` + noteID)

	if err != nil {
		log.Fatal(err)
		return false
	}
	for rows.Next() {

		err = rows.Scan(&setting.OwnerID, &setting.SharedUserID, &setting.Read, &setting.Write)
		if err != nil {
			log.Fatal(err)
			return false
		}
		query := `INSERT INTO SharedSettings (OwnerID, SharedUserID, Read, Write, Name) VALUES ($1, $2, $3, $4, $5)`
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Fatal(err)
			return false
		}
		_, err = stmt.Exec(setting.OwnerID, setting.SharedUserID, setting.Read, setting.Write, setting.Name)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	return true
}

//Logs a user out
func logOut(w http.ResponseWriter, r *http.Request) {
	cookie := checkLoggedIn(r)
	if cookie == nil {
		http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "logged-in",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
		Path:    "/",
	})

	http.Redirect(w, r, "/Users/LogIn", http.StatusSeeOther)
}
