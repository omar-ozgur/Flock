package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
)

func AttendeesIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	attendees := models.GetAttendees()
	j, _ := json.Marshal(attendees)
	w.Write(j)
	fmt.Println("Retrieved attendees")
}

func AttendeesShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	attendee := models.GetAttendee(vars["id"])
	j, _ := json.Marshal(attendee)
	w.Write(j)
	fmt.Println("Retrieved attendee")
}

func AttendeesCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var attendee models.Attendee
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &attendee)
	status := models.CreateAttendee(attendee)
	if status == true {
		fmt.Println("Created new attendee")
	} else {
		fmt.Println("New attendee is not valid")
	}
}

func AttendeesUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var attendee models.Attendee
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &attendee)
	status := models.UpdateAttendee(vars["id"], attendee)
	if status == true {
		fmt.Println("Updated attendee")
	} else {
		fmt.Println("Attendee info is not valid")
	}
}

func AttendeesDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	models.DeleteAttendee(vars["id"])
	fmt.Println("Deleted attendee")
}
