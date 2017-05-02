package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/utilities"
	"os"
	_ "time"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "flock_api"
	DB_HOST     = "localhost"
)

var DB *sql.DB

func InitDB() {
	DBInfo := os.Getenv("DB_INFO")
	if DBInfo == "" {
		DBInfo = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	}
	var err error
	DB, err = sql.Open("postgres", DBInfo)
	utilities.CheckErr(err)

	_, err = DB.Exec("SELECT * FROM users")
	if err != nil {
		_, err = DB.Exec(`CREATE TABLE users (
           id SERIAL,
           first_name text,
           last_name text,
           email text,
           fb_id int,
           time_created timestamp DEFAULT now()
           );`)
		utilities.CheckErr(err)
	}

	_, err = DB.Exec("SELECT * FROM posts")
	if err != nil {
		_, err = DB.Exec(`CREATE TABLE posts (
           id SERIAL,
           title text,
           location text,
           user_id int,
           latitude double precision,
           longitude double precision,
           zip int,           
           time_created timestamp DEFAULT now(),
           time_expires timestamp
           );`)
		utilities.CheckErr(err)
	}

	_, err = DB.Exec("SELECT * FROM attendees")
	if err != nil {
		_, err = DB.Exec(`CREATE TABLE attendees (
           id SERIAL,
           post_id int,
           user_id int
           );`)
		utilities.CheckErr(err)
	}
}
