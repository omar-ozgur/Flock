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

func CreateAttendee(attendee Attendee) (status string, message string, createdAttendee Attendee) {

	// Get attendee fields
	value := reflect.ValueOf(attendee)
	if value.NumField() <= len(attendeeRequiredParams) {
		return "error", "Invalid attendee parameters", Attendee{}
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
			return "error", fmt.Sprintf("Field '%v' is not valid", value.Type().Field(i).Name), Attendee{}
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
	queryStr.WriteString(") returning id;")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	err := db.DB.QueryRow(queryStr.String()).Scan(&attendee.Id)
	if err != nil {
		return "error", "Failed to create new attendee", Attendee{}
	}

	return "success", "New attendee created", attendee
}

func GetPostAttendees(postId string) (status string, message string, retrievedUsers []User) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE post_id='%s';", attendeeTableName, postId)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	rows, err := db.DB.Query(queryStr)
	if err != nil {
		return "error", "Failed to query attendees", nil
	}

	// Print table
	var users []User
	for rows.Next() {
		var attendee Attendee
		err = rows.Scan(&attendee.Id, &attendee.Post_id, &attendee.User_id)
		if err != nil {
			return "error", "Failed to retrieve attendee information", nil
		}
		status, _, retrievedUser := GetUser(fmt.Sprintf("%v", attendee.User_id))
		if status != "success" {
			return "error", "Failed to retrieve attendee user information", nil
		}
		users = append(users, retrievedUser)
	}

	return "success", "Retrieved attendee user information", users
}

func DeleteAttendee(attendee Attendee) (status string, message string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE post_id='%d' AND user_id='%d';", attendeeTableName, attendee.Post_id, attendee.User_id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	_, err := db.DB.Exec(queryStr)
	if err != nil {
		return "error", "Failed to delete attendee"
	}

	return "success", "Deleted attendee"
}
