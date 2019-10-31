package main

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	pass := "123"
	id := 1

	expected := true
	observed := checkPassword(pass, id)

	if observed != expected {
		t.Errorf("Expected true but returned false")
	}
}
