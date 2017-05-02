package models

import (
	"bytes"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"reflect"
)

type Attendee struct {
	Id      int
	Post_id int
	User_id int
}

const attendeeTableName = "attendees"

var attendeeAutoParams = map[string]bool{"Id": true}
var attendeeRequiredParams = map[string]bool{"Post_id": true, "User_id": true}

func CreateAttendee(attendee Attendee) bool {

	// Get attendee fields
	value := reflect.ValueOf(attendee)
	if value.NumField() <= len(attendeeRequiredParams) {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", attendeeTableName))

	// Set present column names
	var first = true
	var values []string
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if attendeeAutoParams[name] {
			continue
		}
		if attendeeRequiredParams[name] && reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
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

func UpdateAttendee(id string, attendee Attendee) bool {

	// Get attendee fields
	value := reflect.ValueOf(attendee)
	if value.NumField() <= 0 {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", attendeeTableName))

	// Set present column names and values
	var first = true
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if attendeeAutoParams[name] {
			continue
		}
		if reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
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

func DeleteAttendee(id string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s", attendeeTableName, id)
	fmt.Println("SQL Query:", queryStr)
	_, err := db.DB.Exec(queryStr)
	utilities.CheckErr(err)
}

func GetAttendee(id string) Attendee {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s", attendeeTableName, id)
	fmt.Println("SQL Query:", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get attendee info
	var attendee Attendee
	err := row.Scan(&attendee.Id, &attendee.Post_id, &attendee.User_id)
	utilities.CheckErr(err)

	return attendee
}

func GetAttendees() []Attendee {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s", attendeeTableName)
	fmt.Println("SQL Query:", queryStr)
	rows, err := db.DB.Query(queryStr)
	utilities.CheckErr(err)

	// Print table
	var attendees []Attendee
	fmt.Printf(" %-5v | %-10v | %-10v\n", "id", "post_id", "user_id")
	for rows.Next() {
		var attendee Attendee
		err = rows.Scan(&attendee.Id, &attendee.Post_id, &attendee.User_id)
		utilities.CheckErr(err)
		attendees = append(attendees, attendee)
		fmt.Printf(" %-5v | %-10v | %-10v\n", attendee.Id, attendee.Post_id, attendee.User_id)
	}

	return attendees
}
