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
	User_id      int       `valid:"required"`
	Latitude     string    `valid:"latitude,required"`
	Longitude    string    `valid:"longitude,required"`
	Zip          int       `valid:"required"`
	Time_created time.Time `valid:"-"`
	Time_expires time.Time `valid:"-"`
}

var eventAutoParams = map[string]bool{"Id": true, "User_id": true, "Time_created": true, "Time_expires": true}
var eventRequiredParams = map[string]bool{"Title": true, "Description": true, "Location": true, "Latitude": true, "Longitude": true, "Zip": true}

func GetEvents() (status string, message string, retrievedEvents []Event) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", utilities.EVENTS_TABLE)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return "error", "Failed to query events", nil
	}

	// Get event info
	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.User_id, &event.Latitude, &event.Longitude, &event.Zip, &event.Time_created, &event.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve event information", nil
		}
		events = append(events, event)
	}

	return "success", "Retrieved events", events
}

func CreateEvent(userId string, event Event) (status string, message string, createdEvent Event) {

	// Get event fields
	value := reflect.ValueOf(event)
	if value.NumField() <= len(eventRequiredParams) {
		return "error", "Invalid event parameters", Event{}
	}

	// Convert user ID to integer
	var err error
	event.User_id, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Event{}
	}

	// Validate user
	_, err = govalidator.ValidateStruct(event)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate event: %s", err.Error()), Event{}
	}

	// Create query string
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
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", value.Field(i).Interface())
		if eventAutoParams[fieldName] {
			continue
		}
		if eventRequiredParams[fieldName] && reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", fieldName), Event{}
		}
		queryStr.WriteString(fmt.Sprintf(", %v", fieldName))
		valuesStr.WriteString(fmt.Sprintf(", $%d", parameterIndex))
		values = append(values, fmt.Sprintf("%v", fieldValue))
		reflections.SetField(&event, value.Type().Field(i).Name, value.Field(i).Interface())
		parameterIndex += 1
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(") VALUES(%s", valuesStr.String()))
	queryStr.WriteString(") returning id;")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Event{}
	}
	err = stmt.QueryRow(values...).Scan(&event.Id)
	if err != nil {
		return "error", "Failed to create new event", Event{}
	}

	// Create attendee
	attendee := Attendee{Event_id: event.Id, User_id: event.User_id}
	fmt.Println(attendee)
	status, _, _ = CreateAttendee(attendee)
	if status == "success" {
		return "success", "New event created", event
	} else {
		return "error", "Failed to add attendee to new event", Event{}
	}
}

func SearchEvents(event Event) (status string, message string, retrievedEvents []Event) {

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", utilities.EVENTS_TABLE))

	// Get event fields
	value := reflect.ValueOf(event)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", nil
	}

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i).Interface()
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(" AND ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)
		parameterIndex += 1
	}

	// Check if any fields are valid
	if len(values) <= 0 {
		return "error", "No valid fields were found", nil
	}

	// Finish and execute query
	queryStr.WriteString(";")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return "error", "Failed to query events", nil
	}

	// Print table
	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.User_id, &event.Latitude, &event.Longitude, &event.Zip, &event.Time_created, &event.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve event information", nil
		}
		events = append(events, event)
	}

	return "success", "Retrieved events", events
}

func GetEvent(id string) (status string, message string, retrievedEvent Event) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", utilities.EVENTS_TABLE)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)
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

func UpdateEvent(id string, event Event) (status string, message string, updatedEvent Event) {

	// Get event fields
	value := reflect.ValueOf(event)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", Event{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", utilities.EVENTS_TABLE))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i).Interface()
		if eventAutoParams[fieldName] {
			continue
		}
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)
		parameterIndex += 1
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id=$%d;", parameterIndex))
	values = append(values, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Event{}
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return "error", "Failed to update event", Event{}
	}

	status, message, retrievedEvent := GetEvent(id)
	if status == "success" {
		return "success", "Updated event", retrievedEvent
	} else {
		return "error", "Failed to retrieve updated event", Event{}
	}
}

func DeleteEvent(id string) (status string, message string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", utilities.EVENTS_TABLE)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return "error", "Failed to delete event"
	}

	return "success", "Deleted event"
}
