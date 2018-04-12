package models

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"reflect"
	"strconv"
	"time"
)

// Event model
type Event struct {
	gorm.Model
	Title        string `valid:"required"`
	Description  string `valid:"required"`
	Location     string `valid:"required"`
	UserID       int
	Attendees    []*User   `gorm:"many2many:attendees;"`
	Latitude     string    `valid:"latitude,required"`
	Longitude    string    `valid:"longitude,required"`
	Zip          int       `valid:"required"`
	Time_expires time.Time `valid:"-"`
}

// eventRequiredParams contains parameters that are required
var eventRequiredParams = map[string]bool{
	"Title":       true,
	"Description": true,
	"Location":    true,
	"Latitude":    true,
	"Longitude":   true,
	"Zip":         true}

// eventModifiableParams contains paramaters that can be modified
var eventModifiableParams = map[string]bool{
	"Title":       true,
	"Description": true,
	"Location":    true,
	"Latitude":    true,
	"Longitude":   true,
	"Zip":         true}

// Initialize events
func InitEvents() {
	Db.AutoMigrate(&Event{})
}

// Get all events
func GetEvents() (status string, message string, retrievedEvents []Event) {
	Db.Find(&retrievedEvents)
	return "success", "Retrieved events", retrievedEvents
}

// checkValid checks if an event is valid
func (event *Event) checkValid() (status string, message string) {
	_, err := govalidator.ValidateStruct(event)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate event: %s", err.Error())
	}
	return "success", ""
}

// Create an event
func CreateEvent(userId string, event Event) (status string, message string, createdEvent Event) {

	// Convert the user ID to an integer
	var err error
	event.UserID, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Event{}
	}

	// Get the user
	status, message, foundUser := GetUser(userId)
	if status != "success" {
		return status, message, Event{}
	}

	// Check if the event is valid
	status, message = event.checkValid()
	if status != "success" {
		return status, message, Event{}
	}

	// Create the event in the database
	Db.Create(&event)

	// Create an attendee
	Db.Model(&event).Association("Attendees").Append(foundUser)

	// Get the created event
	Db.First(&createdEvent, event.ID)
	if status != "success" {
		return "error", "Failed to create the event", Event{}
	}

	return "success", "New event created", createdEvent
}

// Search events
func SearchEvents(parameters map[string]interface{}) (status string, message string, retrievedEvents []Event) {

	Db.Where(parameters).First(&retrievedEvents)

	return "success", "Retrieved events", retrievedEvents
}

// Get an event
func GetEvent(id string) (status string, message string, retrievedEvent Event) {
	Db.First(&retrievedEvent, id)
	return "success", "Retrieved event", retrievedEvent
}

// Update an event
func UpdateEvent(id string, event Event) (status string, message string, updatedEvent Event) {

	// Get event fields
	fields := reflect.ValueOf(event)
	if fields.NumField() <= 0 {
		return "error", "Invalid number of fields", event
	}

	// Check for fields that can't be modified
	parameterIndex := 1
	for i := 0; i < fields.NumField(); i++ {

		// Get field name and value
		fieldName := fields.Type().Field(i).Name

		// Check if field cannot be modified
		if _, ok := eventModifiableParams[fieldName]; !ok {
			return "error", fmt.Sprintf("Field '%s' cannot be modified", fieldName), Event{}
		}

		parameterIndex += 1
	}

	Db.Save(&event)

	// Get the updated event
	Db.First(&updatedEvent, id)
	if status == "success" {
		return "success", "Updated event", updatedEvent
	} else {
		return "error", "Failed to retrieve updated event", Event{}
	}
}

// Delete an event
func DeleteEvent(id string) (status string, message string) {

	// Get event
	status, message, foundEvent := GetEvent(id)
	if status != "success" {
		return status, message
	}

	// Delete associations
	Db.Model(&foundEvent).Association("Attendees").Clear()

	// Delete the event
	Db.Delete(&foundEvent)

	return "success", "Deleted event"
}

func GetAttendees(eventId string) (status string, message string, attendees []User) {

	// Get event
	status, message, foundEvent := GetEvent(eventId)
	if status != "success" {
		return status, message, nil
	}

	// Get attendees
	Db.Model(&foundEvent).Association("Attendees").Find(&attendees)
	if len(attendees) <= 0 {
		return "error", "No attendees found", nil
	}

	return "success", "Found attendees", attendees
}
