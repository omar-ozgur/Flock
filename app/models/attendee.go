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

// Attendee model
type Attendee struct {
	Id       int `valid:"-"`
	Event_id int `valid:"required"`
	User_id  int `valid:"required"`
}

// Parameters that are created automatically
var attendeeAutoParams = map[string]bool{"Id": true}

// Parameters that must be unique
var attendeeUniqueParams = map[string]bool{"Event_id": true, "User_id": true}

// Parameters that are required
var attendeeRequiredParams = map[string]bool{"Event_id": true, "User_id": true}

// Create an attendee
func CreateAttendee(attendee Attendee) (status string, message string, createdAttendee Attendee) {

	// Get attendee fields
	fields := reflect.ValueOf(attendee)
	if fields.NumField() <= len(attendeeRequiredParams) {
		return "error", "Invalid attendee parameters", Attendee{}
	}

	// Validate the attendee
	_, err := govalidator.ValidateStruct(attendee)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate attendee: %s", err.Error()), Attendee{}
	}

	// Check attendee uniqueness
	uniqueMap := make(map[string]interface{})
	for key, _ := range attendeeUniqueParams {

		// Add the field to the unique map
		fieldValue, err := reflections.GetField(&attendee, key)
		if err != nil {
			return "error", fmt.Sprintf("Failed to get field: %s", err.Error()), Attendee{}
		}
		uniqueMap[key] = fieldValue
	}

	// Search for attendees with the same unique values
	status, message, retrievedAttendees := SearchAttendees(uniqueMap, "AND")
	if status != "success" {
		return "error", "Failed to check attendee uniqueness", Attendee{}
	} else if retrievedAttendees != nil {
		return "error", "Attendee already exists", Attendee{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", utilities.ATTENDEES_TABLE))

	// Set present column names
	var fieldsStr, valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < fields.NumField(); i++ {

		// Get the field name and value
		fieldName := fields.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", fields.Field(i).Interface())

		// Skip the field if it is automatically set by the database
		if attendeeAutoParams[fieldName] {
			continue
		}

		// Check if the field is required but empty
		if attendeeRequiredParams[fieldName] && reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", fieldName), Attendee{}
		}

		// Add query delimiters
		if !first {
			fieldsStr.WriteString(", ")
			valuesStr.WriteString(", ")
		} else {
			first = false
		}

		// Add the field name and value to the query
		fieldsStr.WriteString(fieldName)
		valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
		values = append(values, fieldValue)

		parameterIndex += 1
	}

	// Finish the query
	queryStr.WriteString(fmt.Sprintf("%s) VALUES(%s) RETURNING id;", fieldsStr.String(), valuesStr.String()))

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Attendee{}
	}

	// Execute the query
	err = stmt.QueryRow(values...).Scan(&attendee.Id)
	if err != nil {
		return "error", fmt.Sprintf("Failed to create new attendee: %s", err.Error()), Attendee{}
	}

	return "success", "New attendee created", attendee
}

// Get event attendees
func GetEventAttendees(eventId string) (status string, message string, retrievedUsers []User) {

	// Create a query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE event_id=$1;", utilities.ATTENDEES_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", eventId)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}

	// Execute the query
	rows, err := stmt.Query(eventId)
	if err != nil {
		return "error", "Failed to query attendees", nil
	}

	// Get the attendees
	var users []User
	for rows.Next() {

		// Parse the attendee
		var attendee Attendee
		err = rows.Scan(&attendee.Id, &attendee.Event_id, &attendee.User_id)
		if err != nil {
			return "error", "Failed to retrieve attendee information", nil
		}

		// Get the associated user
		status, _, retrievedUser := GetUser(fmt.Sprintf("%v", attendee.User_id))
		if status == "success" {
			users = append(users, retrievedUser)
		}
	}

	return "success", "Retrieved attendee user information", users
}

// Delete an attendee
func DeleteAttendee(attendee Attendee) (status string, message string) {

	// Create the query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE event_id=$1 AND user_id=$2;", utilities.ATTENDEES_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: [%d, %d]", attendee.Event_id, attendee.User_id)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}

	// Execute the query
	_, err = stmt.Exec(attendee.Event_id, attendee.User_id)
	if err != nil {
		return "error", "Failed to delete attendee"
	}

	return "success", "Deleted attendee"
}

// Search for attendees
func SearchAttendees(parameters map[string]interface{}, operator string) (status string, message string, retrievedAttendees []Attendee) {

	// Create a query string
	var queryStr bytes.Buffer

	// Create the query
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", utilities.ATTENDEES_TABLE))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for key, value := range parameters {

		// Add query delimiters
		if !first {
			queryStr.WriteString(fmt.Sprintf(" %s ", operator))
		} else {
			queryStr.WriteString(" ")
			first = false
		}

		// Add the field name and value to the query
		queryStr.WriteString(fmt.Sprintf("%v=$%d", key, parameterIndex))
		values = append(values, value)

		parameterIndex += 1
	}

	// Finish the query
	queryStr.WriteString(";")

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}

	// Execute the query
	rows, err := stmt.Query(values...)
	if err != nil {
		return "error", "Failed to query attendees", nil
	}

	// Get attendees
	var attendees []Attendee
	for rows.Next() {

		// Parse the attendee
		var attendee Attendee
		err = rows.Scan(&attendee.Id, &attendee.Event_id, &attendee.User_id)
		if err != nil {
			return "error", "Failed to retrieve attendee information", nil
		}

		// Add the attendee to the attendees list
		attendees = append(attendees, attendee)
	}

	return "success", "Retrieved attendees", attendees
}
