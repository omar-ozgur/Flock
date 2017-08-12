package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
	"io/ioutil"
	"net/http"
)

// Create a new user
var UsersCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Retrieve body parameters
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)

	// Create user
	status, message, createdUser := models.CreateUser(user)

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"user":    createdUser,
	})
	w.Write(JSON)
})

// Login user
var UsersLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Retrieve body parameters
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)

	// Create user
	status, message, loginToken := models.LoginUser(user)

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"token":   loginToken,
	})
	w.Write(JSON)
})

// Get current user's profile
var UsersProfile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Get user information
	status, message, retrievedUser := models.GetUser(current_user_id)

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"user":    retrievedUser,
	})
	w.Write(JSON)
})

// Get all users
var UsersIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get users
	status, message, retrievedUsers := models.GetUsers()

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"users":   retrievedUsers,
	})
	w.Write(JSON)
})

// Search for users
var UsersSearch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get body parameters
	b, _ := ioutil.ReadAll(r.Body)
	params := make(map[string]interface{})
	json.Unmarshal(b, &params)

	// Search for users
	status, message, retrievedUsers := models.SearchUsers(params, "AND")

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"users":   retrievedUsers,
	})
	w.Write(JSON)
})

// Get a user
var UsersShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get user
	status, message, retrievedUser := models.GetUser(vars["id"])

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"user":    retrievedUser,
	})
	w.Write(JSON)
})

// Update a user
var UsersUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get body and request parameters
	var user models.User
	vars := mux.Vars(r)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Check if the current user has permission to update the specified user
	if vars["id"] != current_user_id {
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "You do not have permission to update this user",
			"user":    models.User{},
		})
		w.Write(JSON)
		return
	}

	// Update the user
	status, message, updatedUser := models.UpdateUser(vars["id"], user)

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"user":    updatedUser,
	})
	w.Write(JSON)
})

// Delete a user
var UsersDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get user claims
	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	// Check if the current user has permission to delete the specified user
	if vars["id"] != current_user_id {
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "You do not have permission to delete this user",
		})
		w.Write(JSON)
		return
	}

	// Delete the user
	status, message := models.DeleteUser(vars["id"])

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
	})
	w.Write(JSON)
})

// Get events a user is going to
var UsersAttendance = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get request parameters
	vars := mux.Vars(r)

	// Get events the user is going to
	status, message, retrievedEvents := models.GetUserAttendance(vars["id"])

	// Return response
	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"events":  retrievedEvents,
	})
	w.Write(JSON)
})
