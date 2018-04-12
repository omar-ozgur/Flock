package models

// CreateAttendance adds the specified user as an attendee for the event
func CreateAttendance(userId string, eventId string) (status string, message string) {

	// Begin the transaction
	tx := Db.Begin()

	// Get the user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(userId); status != "success" {
		return status, message
	}

	// Get the event
	var retrievedEvent Event
	if status, message, retrievedEvent = GetEvent(eventId); status != "success" {
		return status, message
	}

	// Create attendance
	if err := tx.Model(&retrievedUser).Association("Attendances").Append(retrievedEvent).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to create attendance"
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Attendance was successfully recorded"
}

// GetEventAttendees gets attendees for a specific event
func GetEventAttendees(eventId string) (status string, message string, retrievedAttendees []User) {

	// Get event
	var retrievedEvent Event
	if status, message, retrievedEvent = GetEvent(eventId); status != "success" {
		return status, message, nil
	}

	// Get attendees
	if err := Db.Model(&retrievedEvent).Association("Attendees").Find(&retrievedAttendees).Error; err != nil {
		return "error", "Failed to retrieve attendees", nil
	}

	return "success", "Found attendees", retrievedAttendees
}

// GetUserAttendance gets events that a specific user is attending
func GetUserAttendance(userId string) (status string, message string, retrievedEvents []Event) {

	// Get the user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(userId); status != "success" {
		return status, message, nil
	}

	// Find events
	if err := Db.Model(&retrievedUser).Association("Attendances").Find(&retrievedEvents).Error; err != nil {
		return "error", "Failed to retrieve events", nil
	}

	return "success", "Retrieved events", retrievedEvents
}

// DeleteAttendance removes the specified user from the event's attendee list
func DeleteAttendance(userId string, eventId string) (status string, message string) {

	// Begin the transaction
	tx := Db.Begin()

	// Get the user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(userId); status != "success" {
		return status, message
	}

	// Get the event
	var retrievedEvent Event
	if status, message, retrievedEvent = GetEvent(eventId); status != "success" {
		return status, message
	}

	// Delete attendance
	if err := tx.Model(&retrievedUser).Association("Attendances").Delete(retrievedEvent).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to delete attendance"
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Attendance was successfully deleted"
}
