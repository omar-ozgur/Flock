package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/oleiade/reflections.v1"
	"strconv"
	"strings"
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

	// Begin the transaction
	tx := Db.Begin()

	// Convert the user ID to an integer
	var err error
	if event.UserID, err = strconv.Atoi(userId); err != nil {
		return "error", "Invalid user ID", Event{}
	}

	// Get the user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(userId); status != "success" {
		return status, message, Event{}
	}

	// Check if the event is valid
	if status, message = CheckValid(event); status != "success" {
		return status, message, Event{}
	}

	// Create the user in the database
	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to create event", Event{}
	}

	// Create an attendee
	if err := tx.Model(&event).Association("Attendees").Append(retrievedUser).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to create attendee", Event{}
	}

	// Commit the transaction
	tx.Commit()
	return "success", "New event created", event
}

// GetEvent gets a specific event
func GetEvent(id string) (status string, message string, retrievedEvent Event) {

	// Find the event
	if err := Db.First(&retrievedEvent, id).Error; err != nil {
		return "error", "Failed to retrieve event", Event{}
	}

	return "success", "Retrieved event", retrievedEvent
}

// GetEvents gets all events
func GetEvents() (status string, message string, retrievedEvents []Event) {

	// Find the events
	if err := Db.Find(&retrievedEvents).Error; err != nil {
		return "error", "Failed to retrieve events", nil
	}

	return "success", "Retrieved events", retrievedEvents
}

// SearchEvents searches for events
func SearchEvents(params map[string]interface{}) (status string, message string, retrievedEvents []Event) {

	// Modify parameters for sql
	modified := make(map[string]interface{}, len(params))
	for key, value := range params {
		modified[strings.ToLower(key)] = value
	}
	params = modified

	// Search for events
	if err := Db.Where(params).Find(&retrievedEvents).Error; err != nil {
		return "error", "Failed to retrieve events", nil
	}

	return "success", "Retrieved events", retrievedEvents
}

// UpdateEvent updates an event
func UpdateEvent(id string, params map[string]interface{}) (status string, message string, updatedEvent Event) {

	// Begin the transaction
	tx := Db.Begin()

	// Find the existing event
	var retrievedEvent Event
	if status, message, retrievedEvent = GetEvent(id); status != "success" {
		return status, message, Event{}
	}

	// Set changed parameters
	for key, value := range params {
		if err := reflections.SetField(&retrievedEvent, key, value); err != nil {
			return "error", err.Error(), Event{}
		}
	}

	// Check if the event is valid
	if status, message = CheckValid(retrievedEvent); status != "success" {
		return status, message, Event{}
	}

	// Update the event
	if err := tx.Model(&retrievedEvent).Updates(retrievedEvent).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to update the event", Event{}
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Updated event", retrievedEvent
}

// DeleteEvent deletes an event
func DeleteEvent(id string) (status string, message string) {

	// Begin the transaction
	tx := Db.Begin()

	// Find the event
	var retrievedEvent Event
	if status, message, retrievedEvent = GetEvent(id); status != "success" {
		return status, message
	}

	// Delete attendees
	if err := tx.Model(&retrievedEvent).Association("Attendees").Clear().Error; err != nil {
		tx.Rollback()
		return "error", "Failed to delete attendees"
	}

	// Delete the event
	if err := tx.Delete(&retrievedEvent).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to delete the event"
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Deleted event"
}
