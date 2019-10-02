package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func openDB() (db *sql.DB) {
	//Opens database called "EnterpriseNoteApp"
	db, err := sql.Open("portgres", "user=postgres password=password dbname=EnterpriseNoteApp sslmode=diable")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
