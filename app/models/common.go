package models

import (
	"database/sql"
)

var Db *sql.DB

func SetDb(db *sql.DB) {
	Db = db
}

func Init() {
	InitUsers()
	InitEvents()
	InitAttendees()
}
