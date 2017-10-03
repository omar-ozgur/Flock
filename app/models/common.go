package models

import (
	"database/sql"
)

// Db is the database that the models will use
var Db *sql.DB

// SetDb sets the database for the models to use
func SetDb(db *sql.DB) {
	Db = db
}

// Init initializes all models
func Init() {
	InitUsers()
	InitEvents()
	InitAttendees()
}
