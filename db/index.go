package db

import (
	"database/sql"
	"fmt"
	"github.com/omar-ozgur/flock-api/utilities"
	"os"
)

// The database object
var DB *sql.DB

// Initialize the database
func InitDB() {

	// Get database information
	DBInfo := os.Getenv("DB_INFO")
	if DBInfo == "" {
		DBInfo = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
			utilities.DB_USER, utilities.DB_PASSWORD, utilities.DB_NAME, utilities.DB_HOST)
	}

	// Open the database
	DB, err := sql.Open("postgres", DBInfo)
	utilities.CheckErr(err)

	// Create users table if it doesn't exist
	_, err = DB.Exec(fmt.Sprintf("SELECT * FROM %s", utilities.USERS_TABLE))
	if err != nil {
		_, err = DB.Exec(fmt.Sprintf(`CREATE TABLE %s (
           id SERIAL,
           first_name text,
           last_name text,
           email text,
           fb_id text,
           password bytea,
           time_created timestamp DEFAULT now()
           );`, utilities.USERS_TABLE))
		utilities.CheckErr(err)
	}

	// Create events table if it doesn't exist
	_, err = DB.Exec(fmt.Sprintf("SELECT * FROM %s", utilities.EVENTS_TABLE))
	if err != nil {
		_, err = DB.Exec(fmt.Sprintf(`CREATE TABLE %s (
           id SERIAL,
           title text,
           description text,
           location text,
           user_id int,
           latitude text,
           longitude text,
           zip int,           
           time_created timestamp DEFAULT now(),
           time_expires timestamp DEFAULT now()
           );`, utilities.EVENTS_TABLE))
		utilities.CheckErr(err)
	}

	// Create attendees table if it doesn't exist
	_, err = DB.Exec(fmt.Sprintf("SELECT * FROM %s", utilities.ATTENDEES_TABLE))
	if err != nil {
		_, err = DB.Exec(fmt.Sprintf(`CREATE TABLE %s (
           id SERIAL,
           event_id int,
           user_id int
           );`, utilities.ATTENDEES_TABLE))
		utilities.CheckErr(err)
	}
}
