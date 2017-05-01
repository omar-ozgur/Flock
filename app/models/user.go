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

func CreateUser(user User) bool {
	var lastInsertId int
	if user.First_name == "" || user.Last_name == "" || user.Email == "" || user.Fb_id == 0 {
		return false
	}
	err := db.DB.QueryRow("INSERT INTO users(first_name, last_name, email, fb_id, time_created) VALUES($1, $2, $3, $4, $5) returning id;", user.First_name, user.Last_name, user.Email, user.Fb_id, time.Now()).Scan(&lastInsertId)
	utilities.CheckErr(err)
	return true
}

func UpdateUser(id string, user User) bool {
	var err error

	if user.First_name != "" {
		_, err = db.DB.Exec(fmt.Sprintf("update users set first_name='%s' where id=%s", user.First_name, id))
	}
	utilities.CheckErr(err)

	if user.Last_name != "" {
		_, err = db.DB.Exec(fmt.Sprintf("update users set last_name='%s' where id=%s", user.Last_name, id))
	}
	utilities.CheckErr(err)

	if user.Email != "" {
		_, err = db.DB.Exec(fmt.Sprintf("update users set email='%s' where id=%s", user.Email, id))
	}
	utilities.CheckErr(err)

	if user.Fb_id != 0 {
		_, err = db.DB.Exec(fmt.Sprintf("update users set fb_id='%d' where id=%s", user.Fb_id, id))
	}
	utilities.CheckErr(err)

	return true

}

func DeleteUser(id string) {
	_, err := db.DB.Exec(fmt.Sprintf("delete from users where id=%s", id))
	utilities.CheckErr(err)
}

func QueryUsers() []User {
	fmt.Println("# Querying Users")
	rows, err := db.DB.Query("SELECT * FROM users")
	utilities.CheckErr(err)

	var users []User

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
		user := User{Id: id, First_name: first_name, Last_name: last_name, Email: email, Fb_id: fb_id, Time_created: time_created}
		users = append(users, (user))
		fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", id, first_name, last_name, email, fb_id, time_created)
	}

	return users
}
