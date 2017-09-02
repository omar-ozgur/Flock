package utilities

import (
	"fmt"
	"os"
)

// USERS_TABLE is the name of the users table in the database
var USERS_TABLE = "users"

// EVENTS_TABLE is the name of the events table in the database
var EVENTS_TABLE = "events"

// ATTENDEES_TABLE is the name of the attendees table in the database
var ATTENDEES_TABLE = "attendees"

// DB_USER is the username of the database owner
var DB_USER = "postgres"

// DB_PASSWORD is the password of the database owner
var DB_PASSWORD = "postgres"

// DB_NAME is the name of the database
var DB_NAME = "flock_api"

// DB_HOST is the database host
var DB_HOST = "localhost"

// DB_INFO is the string representing the database information
var DB_INFO = func() string {

	// Attempt to retrieve database information from the environment
	dbInfo := os.Getenv("DB_INFO")

	// Use config values if no environment information was found
	if dbInfo == "" {
		dbInfo = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_NAME,
			DB_HOST)
	}

	return dbInfo
}()
