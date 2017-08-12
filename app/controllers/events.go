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

var EventsIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, message, retrievedEvents := models.GetEvents()

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"events":  retrievedEvents,
	})
	w.Write(JSON)
})

var EventsCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	status, message, createdEvent := models.CreateEvent(current_user_id, event)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"event":   createdEvent,
	})
	w.Write(JSON)
})

var EventsSearch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	status, message, retrievedEvents := models.SearchEvents(event)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"events":  retrievedEvents,
	})
	w.Write(JSON)
})

var EventsShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	status, message, retrievedEvent := models.GetEvent(vars["id"])

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"event":   retrievedEvent,
	})
	w.Write(JSON)
})

var EventsUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	status, message, retrievedEvent := models.GetEvent(vars["id"])

	if current_user_id != fmt.Sprintf("%v", retrievedEvent.User_id) {
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "You do not have permission to edit this event",
		})
		w.Write(JSON)
		return
	}

	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	status, message, updatedEvent := models.UpdateEvent(vars["id"], event)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"event":   updatedEvent,
	})
	w.Write(JSON)
})

var EventsDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	status, message, retrievedEvent := models.GetEvent(vars["id"])

	if current_user_id != fmt.Sprintf("%v", retrievedEvent.User_id) {
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "You do not have permission to delete this event",
		})
		w.Write(JSON)
		return
	}

	status, message = models.DeleteEvent(vars["id"])

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
	})
	w.Write(JSON)
})

var EventsAttendees = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var event models.Event
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &event)

	status, message, retrievedAttendees := models.GetEventAttendees(vars["id"])

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":    status,
		"message":   message,
		"attendees": retrievedAttendees,
	})
	w.Write(JSON)
})

var EventsAttend = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	vars := mux.Vars(r)

	eventId, _ := strconv.Atoi(vars["id"])
	userId, _ := strconv.Atoi(current_user_id)

	attendee := models.Attendee{Event_id: eventId, User_id: userId}

	status, message, createdAttendee := models.CreateAttendee(attendee)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":   status,
		"message":  message,
		"attendee": createdAttendee,
	})
	w.Write(JSON)
})

var EventsDeleteAttendance = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	vars := mux.Vars(r)

	eventId, _ := strconv.Atoi(vars["id"])
	userId, _ := strconv.Atoi(current_user_id)

	attendee := models.Attendee{Event_id: eventId, User_id: userId}

	status, message := models.DeleteAttendee(attendee)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
	})
	w.Write(JSON)
})
