package models

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"gopkg.in/oleiade/reflections.v1"
	"reflect"
)

type Attendee struct {
	Id      int `valid:"-"`
	Post_id int `valid:"-"`
	User_id int `valid:"-"`
}

const attendeeTableName = "attendees"

var attendeeAutoParams = map[string]bool{"Id": true}
var attendeeUniqueParams = map[string]bool{"Post_id": true, "User_id": true}
var attendeeRequiredParams = map[string]bool{"Post_id": true, "User_id": true}

func CreateAttendee(attendee Attendee) (status string, message string, createdAttendee Attendee) {

	// Get attendee fields
	value := reflect.ValueOf(attendee)
	if value.NumField() <= len(attendeeRequiredParams) {
		return "error", "Invalid attendee parameters", Attendee{}
	}

	// Validate attendee
	_, err := govalidator.ValidateStruct(attendee)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate attendee: %s", err.Error()), Attendee{}
	}

	// Check attendee uniqueness
	uniqueMap := make(map[string]interface{})
	for key, _ := range attendeeUniqueParams {
		fieldValue, err := reflections.GetField(&attendee, key)
		if err != nil {
			return "error", fmt.Sprintf("Failed to get field: %s", err.Error()), Attendee{}
		}
		uniqueMap[key] = fieldValue
	}
	status, message, retrievedAttendees := SearchAttendees(uniqueMap, "AND")
	if status != "success" {
		return "error", "Failed to check attendee uniqueness", Attendee{}
	} else if retrievedAttendees != nil {
		return "error", "Attendee already exists", Attendee{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", attendeeTableName))

	// Set present column names
	var fieldsStr, valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", value.Field(i).Interface())
		if attendeeAutoParams[fieldName] {
			continue
		}
		if attendeeRequiredParams[fieldName] && reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", fieldName), Attendee{}
		}
		if !first {
			fieldsStr.WriteString(", ")
			valuesStr.WriteString(", ")
		} else {
			first = false
		}
		fieldsStr.WriteString(fieldName)
		valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
		parameterIndex += 1
		values = append(values, fieldValue)
	}

	// Finish and prepare query
	queryStr.WriteString(fmt.Sprintf("%s) VALUES(%s) RETURNING id;", fieldsStr.String(), valuesStr.String()))
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Attendee{}
	}

	// Execute query
	err = stmt.QueryRow(values...).Scan(&attendee.Id)
	if err != nil {
		return "error", fmt.Sprintf("Failed to create new attendee: %s", err.Error()), Attendee{}
	}

	return "success", "New attendee created", attendee
}

func GetPostAttendees(postId string) (status string, message string, retrievedUsers []User) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE post_id=$1;", attendeeTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", postId)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query(postId)
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
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE post_id=$1 AND user_id=$2;", attendeeTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: [%d, %d]", attendee.Post_id, attendee.User_id)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}
	_, err = stmt.Exec(attendee.Post_id, attendee.User_id)
	if err != nil {
		return "error", "Failed to delete attendee"
	}

	return "success", "Deleted attendee"
}

func SearchAttendees(parameters map[string]interface{}, operator string) (status string, message string, retrievedAttendees []Attendee) {

	// Create query string
	var queryStr bytes.Buffer

	// Create and execute query
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", attendeeTableName))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for key, value := range parameters {
		if !first {
			queryStr.WriteString(fmt.Sprintf(" %s ", operator))
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v=$%d", key, parameterIndex))
		values = append(values, value)
		parameterIndex += 1
	}

	queryStr.WriteString(";")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return "error", "Failed to query attendees", nil
	}

	// Print table
	var attendees []Attendee
	for rows.Next() {
		var attendee Attendee
		err = rows.Scan(&attendee.Id, &attendee.Post_id, &attendee.User_id)
		if err != nil {
			return "error", "Failed to retrieve attendee information", nil
		}
		attendees = append(attendees, attendee)
	}

	return "success", "Retrieved attendees", attendees
}
