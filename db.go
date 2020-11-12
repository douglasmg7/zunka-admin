package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Init zunka db.
func initZunkaDB() {
	// Start data base.
	dbZunka, err = sql.Open("sqlite3", dbZunkaFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Error open db %v: %v", dbZunka, err))
	}
	err = dbZunka.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Close zunka db.
func closeZunkaDB() {
	dbZunka.Close()
}

// Init aldo db.
func initAldoDB() {
	dbAldo = sqlx.MustConnect("sqlite3", dbAldoFile)
}

// Close aldo db.
func closeAldoDB() {
	dbAldo.Close()
}

// Init allnations db.
func initAllnationsDB() {
	dbAllnations = sqlx.MustConnect("sqlite3", dbAllnationsFile)
}

// Close allnations db.
func closeAllnationsDB() {
	dbAllnations.Close()
}
