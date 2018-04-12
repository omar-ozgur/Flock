package models

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Attendee model
type Attendee struct {
	gorm.Model
	Event_id int `valid:"required"`
	User_id  int `valid:"required"`
}

// attendeeRequiredParams contains parameters that are required
var attendeeRequiredParams = map[string]bool{"Event_id": true, "User_id": true}

// attendeeModifiableParams contains paramaters that can be modified
var attendeeModifiableParams = map[string]bool{"Event_id": true, "User_id": true}

// Initialize attendees
func InitAttendees() {
	Db.AutoMigrate(&Attendee{})
}

// checkValid checks if an event is valid
func (attendee *Attendee) checkValid() (status string, message string) {
	_, err := govalidator.ValidateStruct(attendee)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate attendee: %s", err.Error())
	}
	return "success", ""
}

// Create an attendee
func CreateAttendee(attendee Attendee) (status string, message string, createdAttendee Attendee) {

	// Check if the attendee is valid
	status, message = attendee.checkValid()
	if status != "success" {
		return status, message, Attendee{}
	}

	// Create the attendee in the database
	Db.Create(&attendee)

	// Get the created attendee
	Db.First(&createdAttendee, attendee.ID)
	if status != "success" {
		return "error", "Failed to create the attendee", Attendee{}
	}

	return "success", "New attendee created", createdAttendee
}

// Get event attendees
func GetEventAttendees(eventId string) (status string, message string, retrievedUsers []User) {

	var retrievedAttendees []Attendee
	Db.Where("event_id = ?", eventId).Find(&retrievedAttendees)

	// Get the attendees
	for _, attendee := range retrievedAttendees {

		// Get the associated user
		status, _, retrievedUser := GetUser(fmt.Sprintf("%v", attendee.User_id))
		if status == "success" {
			retrievedUsers = append(retrievedUsers, retrievedUser)
		}
	}

	return "success", "Retrieved attendee user information", retrievedUsers
}

// Delete an attendee
func DeleteAttendee(attendee Attendee) (status string, message string) {
	Db.Delete(&attendee)
	return "success", "Deleted attendee"
}

// Search for attendees
func SearchAttendees(parameters map[string]interface{}) (status string, message string, retrievedAttendees []Attendee) {

	Db.Where(parameters).First(&retrievedAttendees)

	return "success", "Retrieved attendees", retrievedAttendees
}
