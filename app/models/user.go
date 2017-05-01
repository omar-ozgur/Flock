package models

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"time"
)

type User struct {
	Id           int
	First_name   string
	Last_name    string
	Email        string
	Fb_id        int
	Time_created time.Time
}

func SaveUser(user User) {
	var lastInsertId int
	err := db.DB.QueryRow("INSERT INTO users(first_name, last_name, email, fb_id, time_created) VALUES($1, $2, $3, $4, $5) returning id;", user.First_name, user.Last_name, user.Email, user.Fb_id, time.Now()).Scan(&lastInsertId)
	utilities.CheckErr(err)
	fmt.Println("Inserted new user")
}

func QueryUsers() {
	fmt.Println("# Querying Users")
	rows, err := db.DB.Query("SELECT * FROM users")
	utilities.CheckErr(err)

	fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", "id", "first_name", "last_name", "email", "fb_id", "time_created")
	for rows.Next() {
		var id int
		var first_name string
		var last_name string
		var email string
		var fb_id int
		var time_created time.Time
		err = rows.Scan(&id, &first_name, &last_name, &email, &fb_id, &time_created)
		utilities.CheckErr(err)
		fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", id, first_name, last_name, email, fb_id, time_created)
	}
}
