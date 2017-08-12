package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Get all events
var EventsIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
})

// Create an event
var EventsCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get body parameters
	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	// Create an event
	status, message, createdEvent := models.CreateEvent(current_user_id, event)

	// Return a response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"event":   createdEvent,
	})
	w.Write(JSON)
})

// Search for events
var EventsSearch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get body parameters
	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	// Search for events
	status, message, retrievedEvents := models.SearchEvents(event)

	// Return a response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"events":  retrievedEvents,
	})
	w.Write(JSON)
})

// Get an event
var EventsShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
})

// Update an event
var EventsUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get an event
	status, message, retrievedEvent := models.GetEvent(vars["id"])

	// Check if the user has sufficient permissions to update the event
	if current_user_id != fmt.Sprintf("%v", retrievedEvent.User_id) {
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "You do not have permission to edit this event",
		})
		w.Write(JSON)
		return
	}

	// Get body parameters
	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	// Update the event
	status, message, updatedEvent := models.UpdateEvent(vars["id"], event)

	// Return a response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"event":   updatedEvent,
	})
	w.Write(JSON)
})

// Delete an event
var EventsDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get events
	status, message, retrievedEvent := models.GetEvent(vars["id"])

	// Check if the user has sufficient permissions to delete the event
	if current_user_id != fmt.Sprintf("%v", retrievedEvent.User_id) {
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
})

// Get event attendees
var EventsAttendees = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get body parameters
	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	// Get event attendees
	status, message, retrievedAttendees := models.GetEventAttendees(vars["id"])

	// Return a response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":    status,
		"message":   message,
		"attendees": retrievedAttendees,
	})
	w.Write(JSON)
})

// Attend an event
var EventsAttend = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get request parameters
	vars := mux.Vars(r)

	// Get the event Id
	eventId, _ := strconv.Atoi(vars["id"])

	// Get the user Id
	userId, _ := strconv.Atoi(current_user_id)

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
})

// Delete an event attendee
var EventsDeleteAttendance = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get request parameters
	vars := mux.Vars(r)

	// Get the event Id
	eventId, _ := strconv.Atoi(vars["id"])

	// Get the user Id
	userId, _ := strconv.Atoi(current_user_id)

	// Delete the attendee
	attendee := models.Attendee{Event_id: eventId, User_id: userId}
	status, message := models.DeleteAttendee(attendee)

	// Return a response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
	})
	w.Write(JSON)
})
