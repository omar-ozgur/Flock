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
	"strconv"
	"time"
)

// Event model
type Event struct {
	Id           int       `valid:"-"`
	Title        string    `valid:"required"`
	Description  string    `valid:"required"`
	Location     string    `valid:"required"`
	User_id      int       `valid:"-"`
	Latitude     string    `valid:"latitude,required"`
	Longitude    string    `valid:"longitude,required"`
	Zip          int       `valid:"required"`
	Time_created time.Time `valid:"-"`
	Time_expires time.Time `valid:"-"`
}

// Parameters that are created automatically
var eventAutoParams = map[string]bool{"Id": true, "User_id": true, "Time_created": true, "Time_expires": true}

// Parameters that are required
var eventRequiredParams = map[string]bool{"Title": true, "Description": true, "Location": true, "Latitude": true, "Longitude": true, "Zip": true}

// Get all events
func GetEvents() (status string, message string, retrievedEvents []Event) {

	// Create a query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", utilities.EVENTS_TABLE)

	// Log the query
	utilities.Sugar.Infof("SQL Query: %s", queryStr)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		return "error", "Failed to query events", nil
	}

	// Get event info
	var events []Event
	for rows.Next() {

		// Parse the event
		var event Event
		err = rows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.User_id, &event.Latitude, &event.Longitude, &event.Zip, &event.Time_created, &event.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve event information", nil
		}

		// Append the event to the events list
		events = append(events, event)
	}

	return "success", "Retrieved events", events
}

// Create an event
func CreateEvent(userId string, event Event) (status string, message string, createdEvent Event) {

	// Get event fields
	fields := reflect.ValueOf(event)

	// Convert the user ID to an integer
	var err error
	event.User_id, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Event{}
	}

	// Validate the event
	_, err = govalidator.ValidateStruct(event)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate event: %s", err.Error()), Event{}
	}

	// Create a query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", utilities.EVENTS_TABLE))

	// Set present column names
	var valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	queryStr.WriteString("User_id")
	valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
	parameterIndex += 1
	values = append(values, userId)
	for i := 0; i < fields.NumField(); i++ {

		// Get the field name and value
		fieldName := fields.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", fields.Field(i).Interface())

		// Skip the field if it is automatically set by the database
		if eventAutoParams[fieldName] {
			continue
		}

		// Check if the field is empty even though it's required
		if eventRequiredParams[fieldName] && reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", fieldName), Event{}
		}

		// Add names and values to the query
		queryStr.WriteString(fmt.Sprintf(", %v", fieldName))
		valuesStr.WriteString(fmt.Sprintf(", $%d", parameterIndex))
		values = append(values, fmt.Sprintf("%v", fieldValue))

		// Set the event's field
		reflections.SetField(&event, fields.Type().Field(i).Name, fields.Field(i).Interface())

		parameterIndex += 1
	}

	// Finish the query
	queryStr.WriteString(fmt.Sprintf(") VALUES(%s", valuesStr.String()))
	queryStr.WriteString(") returning id;")

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Event{}
	}

	// Execute the query
	err = stmt.QueryRow(values...).Scan(&event.Id)
	if err != nil {
		return "error", "Failed to create new event", Event{}
	}

	// Create an attendee
	attendee := Attendee{Event_id: event.Id, User_id: event.User_id}
	fmt.Println(attendee)
	status, _, _ = CreateAttendee(attendee)
	if status == "success" {
		return "success", "New event created", event
	} else {
		return "error", "Failed to add attendee to new event", Event{}
	}
}

// Search events
func SearchEvents(event Event) (status string, message string, retrievedEvents []Event) {

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", utilities.EVENTS_TABLE))

	// Get event fields
	fields := reflect.ValueOf(event)
	if fields.NumField() <= 0 {
		return "error", "Invalid number of fields", nil
	}

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < fields.NumField(); i++ {

		// Get the field
		fieldName := fields.Type().Field(i).Name
		fieldValue := fields.Field(i).Interface()

		// Skip the field if it is empty
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}

		// Add query delimiters
		if !first {
			queryStr.WriteString(" AND ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}

		// Add the field name and value to the query
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)

		parameterIndex += 1
	}

	// Check if any fields are valid
	if len(values) <= 0 {
		return "error", "No valid fields were found", nil
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
		return "error", "Failed to query events", nil
	}

	// Print table
	var events []Event
	for rows.Next() {

		// Parse the event
		var event Event
		err = rows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.User_id, &event.Latitude, &event.Longitude, &event.Zip, &event.Time_created, &event.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve event information", nil
		}

		// Append the event to the events list
		events = append(events, event)
	}

	return "success", "Retrieved events", events
}

// Get an event
func GetEvent(id string) (status string, message string, retrievedEvent Event) {

	// Create a query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", utilities.EVENTS_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Event{}
	}
	row := stmt.QueryRow(id)

	// Get event info
	var event Event
	err = row.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.User_id, &event.Latitude, &event.Longitude, &event.Zip, &event.Time_created, &event.Time_expires)
	if err != nil {
		return "error", "Failed to retrieve event information", Event{}
	}

	return "success", "Retrieved event", event
}

// Update an event
func UpdateEvent(id string, event Event) (status string, message string, updatedEvent Event) {

	// Get event fields
	fields := reflect.ValueOf(event)
	if fields.NumField() <= 0 {
		return "error", "Invalid number of fields", Event{}
	}

	// Create a query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", utilities.EVENTS_TABLE))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < fields.NumField(); i++ {

		// Get the field name and value
		fieldName := fields.Type().Field(i).Name
		fieldValue := fields.Field(i).Interface()

		// Skip the field if it is automatically set by the database
		if eventAutoParams[fieldName] {
			continue
		}

		// Skip the field if it is empty
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}

		// Add query delimiters
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}

		// Add the field name and value to the query
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)

		parameterIndex += 1
	}

	// Finish the query
	queryStr.WriteString(fmt.Sprintf(" WHERE id=$%d;", parameterIndex))
	values = append(values, id)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Event{}
	}

	// Execute the query
	_, err = stmt.Exec(values...)
	if err != nil {
		return "error", "Failed to update event", Event{}
	}

	// Get the updated event
	status, message, retrievedEvent := GetEvent(id)
	if status == "success" {
		return "success", "Updated event", retrievedEvent
	} else {
		return "error", "Failed to retrieve updated event", Event{}
	}
}

// Delete an event
func DeleteEvent(id string) (status string, message string) {

	// Create a query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", utilities.EVENTS_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}

	// Execute the query
	_, err = stmt.Exec(id)
	if err != nil {
		return "error", "Failed to delete event"
	}

	// Delete the event's attendees
	attendeeParams := make(map[string]interface{})
	attendeeParams["Event_id"] = id
	status, _, retrievedAttendees := SearchAttendees(attendeeParams, "AND")
	if status == "success" {
		for _, attendee := range retrievedAttendees {
			DeleteAttendee(attendee)
		}
	}

	return "success", "Deleted event"
}
