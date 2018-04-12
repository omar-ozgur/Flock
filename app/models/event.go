package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/oleiade/reflections.v1"
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

// InitEvents initializes events
func InitEvents() {
	Db.AutoMigrate(&Event{})
}

// CreateEvent creates an event
func CreateEvent(event Event, userId string) (status string, message string, createdEvent Event) {

	// Convert the user ID to an integer
	var err error
	event.UserID, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Event{}
	}

	// Get the user
	status, message, retrievedUser := GetUser(userId)
	if status != "success" {
		return status, message, Event{}
	}

	// Check if the event is valid
	status, message = CheckValid(event)
	if status != "success" {
		return status, message, Event{}
	}

	// Create the user in the database
	Db.Create(&event)
	createdEvent = event
	if createdEvent.ID == 0 {
		return "error", "Failed to create event", Event{}
	}

	// Create an attendee
	Db.Model(&createdEvent).Association("Attendees").Append(retrievedUser)

	return "success", "New event created", createdEvent
}

// GetEvent gets a specific event
func GetEvent(id string) (status string, message string, retrievedEvent Event) {

	// Find the event
	Db.First(&retrievedEvent, id)
	if retrievedEvent.ID == 0 {
		return "error", "Failed to retrieve event", Event{}
	}

	return "success", "Retrieved event", retrievedEvent
}

// GetEvents gets all events
func GetEvents() (status string, message string, retrievedEvents []Event) {

	// Find the events
	Db.Find(&retrievedEvents)
	if len(retrievedEvents) <= 0 {
		return "error", "Failed to retrieve events", nil
	}

	return "success", "Retrieved events", retrievedEvents
}

// SearchEvents searches for events
func SearchEvents(params map[string]interface{}) (status string, message string, retrievedEvents []Event) {

	// Search for events
	Db.Where(params).First(&retrievedEvents)
	if len(retrievedEvents) <= 0 {
		return "error", "Failed to retrieve events", nil
	}

	return "success", "Retrieved events", retrievedEvents
}

// UpdateEvent updates an event
func UpdateEvent(id string, params map[string]interface{}) (status string, message string, updatedEvent Event) {

	// Find the existing event
	status, message, retrievedEvent := GetEvent(id)
	if status != "success" {
		return status, message, Event{}
	}

	// Set changed parameters
	for key, value := range params {

		// Set updated field values
		err := reflections.SetField(&retrievedEvent, key, value)
		if err != nil {
			return "error", err.Error(), Event{}
		}
	}

	// Check if the event is valid
	status, message = CheckValid(retrievedEvent)
	if status != "success" {
		return status, message, Event{}
	}

	// Update the event
	Db.Model(&retrievedEvent).Updates(retrievedEvent)

	// Get the updated event
	Db.First(&updatedEvent, id)

	return "success", "Updated event", updatedEvent
}

// DeleteEvent deletes an event
func DeleteEvent(id string) (status string, message string) {

	// Find the event
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

// GetAttendees gets attendees for a specific event
func GetAttendees(eventId string) (status string, message string, retrievedAttendees []User) {

	// Get event
	status, message, foundEvent := GetEvent(eventId)
	if status != "success" {
		return status, message, nil
	}

	// Get attendees
	Db.Model(&foundEvent).Association("Attendees").Find(&retrievedAttendees)

	return "success", "Found attendees", retrievedAttendees
}
