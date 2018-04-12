package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

// parseEvent parses an event from the body of a request
func parseEvent(r *http.Request) models.Event {
	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)
	return event
}

// EventsIndex gets information for all events
var EventsIndex = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get events
		status, message, retrievedEvents := models.GetEvents()

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"events":  retrievedEvents,
		})
		w.Write(JSON)
	},
)

// EventsCreate creates a new event
var EventsCreate = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Parse the event from the body
		event := parseEvent(r)

		// Create an event
		status, message, createdEvent := models.CreateEvent(
			currentUserId,
			event)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"event":   createdEvent,
		})
		w.Write(JSON)
	},
)

// EventsSearch searches for events
var EventsSearch = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get body parameters
		b, _ := ioutil.ReadAll(r.Body)
		params := make(map[string]interface{})
		json.Unmarshal(b, &params)

		// Search for events
		status, message, retrievedEvents := models.SearchEvents(params)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"events":  retrievedEvents,
		})
		w.Write(JSON)
	},
)

// EventsShow retrieves information for a specific event
var EventsShow = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get an event
		status, message, retrievedEvent := models.GetEvent(vars["id"])

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"event":   retrievedEvent,
		})
		w.Write(JSON)
	},
)

// EventsUpdate updates a specific event
var EventsUpdate = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get an event
		status, message, retrievedEvent := models.GetEvent(vars["id"])

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Check if the user has sufficient permissions to update the event
		if currentUserId != fmt.Sprintf("%v", retrievedEvent.User_id) {
			JSON, _ := json.Marshal(map[string]interface{}{
				"status":  "error",
				"message": "You do not have permission to edit this event",
			})
			w.Write(JSON)
			return
		}

		// Parse the event from the body
		event := parseEvent(r)

		// Update the event
		status, message, updatedEvent := models.UpdateEvent(vars["id"], event)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"event":   updatedEvent,
		})
		w.Write(JSON)
	},
)

// EventsDelete deletes an event
var EventsDelete = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get events
		status, message, retrievedEvent := models.GetEvent(vars["id"])

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Check if the user has sufficient permissions to delete the event
		if currentUserId != fmt.Sprintf("%v", retrievedEvent.User_id) {
			JSON, _ := json.Marshal(map[string]interface{}{
				"status":  "error",
				"message": "You do not have permission to delete this event",
			})
			w.Write(JSON)
			return
		}

		// Delete the event
		status, message = models.DeleteEvent(vars["id"])

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
		})
		w.Write(JSON)
	},
)

// EventsAttendees retrieves the attendees of an event
var EventsAttendees = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get event attendees
		status, message, retrievedAttendees := models.GetEventAttendees(
			vars["id"])

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":    status,
			"message":   message,
			"attendees": retrievedAttendees,
		})
		w.Write(JSON)
	},
)

// EventsAttend causes the current user to join a specific event
var EventsAttend = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Get request parameters
		vars := mux.Vars(r)

		// Get the event Id
		eventId, _ := strconv.Atoi(vars["id"])

		// Get the user Id
		userId, _ := strconv.Atoi(currentUserId)

		// Create an attendee
		attendee := models.Attendee{Event_id: eventId, User_id: userId}
		status, message, createdAttendee := models.CreateAttendee(attendee)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":   status,
			"message":  message,
			"attendee": createdAttendee,
		})
		w.Write(JSON)
	},
)

// EventsDeleteAttendance removes a user from a specific event's attendee list
var EventsDeleteAttendance = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Get request parameters
		vars := mux.Vars(r)

		// Get the event Id
		eventId, _ := strconv.Atoi(vars["id"])

		// Get the user Id
		userId, _ := strconv.Atoi(currentUserId)

		// Delete the attendee
		attendee := models.Attendee{Event_id: eventId, User_id: userId}
		status, message := models.DeleteAttendee(attendee)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
		})
		w.Write(JSON)
	},
)
