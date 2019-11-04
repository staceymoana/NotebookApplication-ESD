package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var w http.ResponseWriter
var r *http.Request

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	setupDB()
	//os.Exit(m.Run())

}

func TestDatabase(t *testing.T) {
	var user User
	user.UserID = 100
	user.GivenName = "John"
	user.FamilyName = "Snow"
	user.Password = "hello123"

	var note Note
	note.NoteID = 100
	note.UserID = 100
	note.Title = "my note"
	note.Contents = "hi this is a note"
	note.DateCreated = time.Now()
	note.DateUpdated = time.Now()

	db := setupDB()

	if db != nil {
		assert.Equal(t, "Users returned", getUsers(w, r), "Should return 'Users returned'")
	}

}

func TestCheckPassword(t *testing.T) {
	pass := "123"
	id := 1

	expected := true
	observed := checkPassword(pass, id)

	if observed != expected {
		t.Errorf("Expected true but returned false")
	}
}

func TestIsOwner(t *testing.T) {
	result := isOwner(w, r)

	if result == false {
		t.Errorf("Is not owner")
	}
}

func TestSetupDB(t *testing.T) {
	result := setupDB()

	if result == nil {
		t.Errorf("Database is nil")
	}
}
