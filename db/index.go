package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/utilities"
	_ "time"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "flock-api"
	DB_HOST     = "localhost"
)

var DB *sql.DB

func InitDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	var err error
	DB, err = sql.Open("postgres", dbinfo)
	utilities.CheckErr(err)

	_, err = DB.Exec("SELECT * FROM users")
	if err != nil {
		_, err = DB.Exec(`CREATE TABLE users (
           id SERIAL,
           first_name text,
           last_name text,
           email text,
           fb_id int,
           time_created timestamp
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
           time_created timestamp,
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

	// var lastInsertId int
	// err = DB.QueryRow("INSERT INTO users(first_name, last_name, email, fb_id, time_created) VALUES($1, $2, $3, $4, $5) returning id;", "Omar", "Ozgur", "oozgur217@gmail.com", "1", time.Now()).Scan(&lastInsertId)
	// utilities.CheckErr(err)
	// fmt.Println("last inserted id =", lastInsertId)

	// fmt.Println("# Updating")
	// stmt, err := DB.Prepare("update userinfo set username=$1 where uid=$2")
	// utilities.CheckErr(err)

	// res, err := stmt.Exec("astaxieupdate", lastInsertId)
	// utilities.CheckErr(err)

	// affect, err := res.RowsAffected()
	// utilities.CheckErr(err)

	// fmt.Println(affect, "rows changed")

	// fmt.Println("# Querying")
	// rows, err := DB.Query("SELECT * FROM users")
	// utilities.CheckErr(err)

	// fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", "id", "first_name", "last_name", "email", "fb_id", "time_created")
	// for rows.Next() {
	// 	var id int
	// 	var first_name string
	// 	var last_name string
	// 	var email string
	// 	var fb_id int
	// 	var time_created time.Time
	// 	err = rows.Scan(&id, &first_name, &last_name, &email, &fb_id, &time_created)
	// 	utilities.CheckErr(err)
	// 	fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", id, first_name, last_name, email, fb_id, time_created)
	// }

	// fmt.Println("# Deleting")
	// stmt, err = DB.Prepare("delete from userinfo where uid=$1")
	// utilities.CheckErr(err)

	// res, err = stmt.Exec(lastInsertId)
	// utilities.CheckErr(err)

	// affect, err = res.RowsAffected()
	// utilities.CheckErr(err)

	// fmt.Println(affect, "rows changed")
}
