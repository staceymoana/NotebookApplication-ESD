package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var w http.ResponseWriter
var r *http.Request

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	setupDB()
	os.Exit(m.Run())

}

func TestDatabase(t *testing.T) {

	db := setupDB()

	if assert.NotNil(t, db) {
		//getUsersSQL() gets existing users from the database and returns them
		assert.NotNil(t, getUsersSQL(), "getUsersSQL() should return a list of users")
		//getUserNotesSQL() gets notes from the given user using their UserID and returns them
		userNotes := getUserNotesSQL("1")
		assert.NotEmpty(t, userNotes, "getUserNotesSQL() should not be empty")
		//updateNoteInsertSQL() updates a note based on input and returns true if sucessful
		assert.True(t, updateNoteInsertSQL("Updated title", "Updated contents", "1"), "updateNoteInsertSQL() should return true")
		//isOwnerSQL() checks if logged in user is the owner of a note and returns the ownerID
		ownerID := isOwnerSQL("1", "1")
		assert.NotZero(t, ownerID, "isOwnerSQL() should not return zero")
		//deleteNoteSQL() deletes a note based on a given NoteID returns true if sucessful
		assert.True(t, deleteNoteSQL("2"), "deleteNoteSQL() should return true")
		//createUserSQL() creates a new user based on input and returns the user
		newUser := createUserSQL("New", "User", "password")
		assert.NotNil(t, newUser, "createUserSQL() should return a user")
		//searchSQL() searches a note based on input on a given NoteID, returns array of notes containing input
		searchedNotes := searchSQL("content", "1")
		assert.NotEmpty(t, searchedNotes, "searchedNotes() should not be empty")
		//analyseNoteSQL() analyses a note based on input on a given NoteID and returns a count
		newAnalyseNote := analyseNoteSQL("content", "1")
		assert.NotZero(t, newAnalyseNote, "analyseNoteSQL() should not return zero")
		//shareNoteSQL() shares a note based on input and returns true if sucessful
		assert.True(t, shareNoteSQL("1", "on", "on", "4"), "shareNoteSQL() should return true")
		//accessSQL() gets existing access and returns them
		newAccess := accessSQL("1")
		assert.NotEmpty(t, newAccess, "accessSQL() should not be empty")
		//editAccessSQL() edits an access setting and returns true if successful
		assert.True(t, editAccessSQL("on", "on", "4"), "editAccessSQL() should return true")
		//saveSharedSettingOnNoteSQL() saves shared settings based on input and NoteID and returns true if successful
		assert.True(t, saveSharedSettingOnNoteSQL("test", "400"), "saveSharedSettingOnNoteSQL() should return true")
		//createNoteSelectSQL() returns the SharedSettings based on logged in user and returns them as array
		settings := createNoteSelectSQL("1")
		assert.NotEmpty(t, settings, "createNoteSelectSQL() should not be empty")
		//createNoteInsertSQL() creates a new note and returns true if successful
		assert.True(t, createNoteInsertSQL("1", "test title", "test contents", "none"), "createNoteInsertSQL() should return true")
		//updateNoteSelectSQL() updates a note based on a NoteID
		writevalue, id := updateNoteSelectSQL("1")
		//returns true if successfully updated
		assert.True(t, writevalue, "updateNoteSelectSQL() writevalue should return true")
		//check if id is not empty
		assert.NotEqual(t, "", id, "updateNoteSelectSQL() id should not be empty")
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
