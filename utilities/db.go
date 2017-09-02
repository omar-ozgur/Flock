package utilities

import (
	"fmt"
	"os"
)

// Database configuration values
var DB_INFO = getDbInfo()
var DB_USER = "postgres"
var DB_PASSWORD = "postgres"
var DB_NAME = "flock_api"
var DB_HOST = "localhost"

// Table names
var USERS_TABLE = "users"
var EVENTS_TABLE = "events"
var ATTENDEES_TABLE = "attendees"

func getDbInfo() string {
	dbInfo := os.Getenv("DB_INFO")
	if dbInfo == "" {
		dbInfo = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_NAME,
			DB_HOST)
	}
	return dbInfo
}
