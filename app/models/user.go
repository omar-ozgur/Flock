package models

import (
	"bytes"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"reflect"
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

const tableName = "users"

var auto = map[string]bool{"Id": true, "Time_created": true}
var required = map[string]bool{"First_name": true, "Last_name": true, "Email": true, "Fb_id": true}

func CreateUser(user User) bool {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= len(required) {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))

	// Set present column names
	var first = true
	var values []string
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if auto[name] {
			continue
		}
		if required[name] && reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			return false
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v", value.Type().Field(i).Name))
		values = append(values, fmt.Sprintf("%v", value.Field(i).Interface()))
	}

	// Set present column values
	queryStr.WriteString(") VALUES(")
	first = true
	for i := 0; i < len(values); i++ {
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("'%v'", values[i]))
	}

	// Finish and execute query
	queryStr.WriteString(")")
	fmt.Println("SQL Query:", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	utilities.CheckErr(err)

	return true
}

func UpdateUser(id string, user User) bool {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= 0 {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", tableName))

	// Set present column names and values
	var first = true
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if auto[name] {
			continue
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		if reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			return false
		}
		queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, value.Field(i).Interface()))
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id='%s'", id))
	fmt.Println("SQL Query:", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	utilities.CheckErr(err)

	return true
}

func DeleteUser(id string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s", tableName, id)
	fmt.Println("SQL Query:", queryStr)
	_, err := db.DB.Exec(queryStr)
	utilities.CheckErr(err)
}

func GetUser(id string) User {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s", tableName, id)
	fmt.Println("SQL Query:", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get user info
	var user User
	err := row.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Time_created)
	utilities.CheckErr(err)

	return user
}

func GetUsers() []User {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s", tableName)
	fmt.Println("SQL Query:", queryStr)
	rows, err := db.DB.Query(queryStr)
	utilities.CheckErr(err)

	// Print table
	var users []User
	fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", "id", "first_name", "last_name", "email", "fb_id", "time_created")
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Time_created)
		utilities.CheckErr(err)
		users = append(users, user)
		fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v\n", user.Id, user.First_name, user.Last_name, user.Email, user.Fb_id, user.Time_created)
	}

	return users
}
