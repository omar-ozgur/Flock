package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
	"io/ioutil"
	"net/http"
)

// parseUser parses a user from the body of a request
func parseUser(r *http.Request) models.User {
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	return user
}

// parseCurrentUser parses the current user based on authorization information
// Returns the current user's ID
func parseCurrentUser(r *http.Request) string {
	claims := utilities.GetClaims(
		r.Header.Get("Authorization")[len("Bearer "):])
	return fmt.Sprintf("%v", claims["user_id"])
}

// UsersCreate creates a new user
var UsersCreate = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Parse the user from the body
		user := parseUser(r)

		// Create the user
		status, message, createdUser := models.CreateUser(user)

		// Return the response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"user":    createdUser,
		})
		w.Write(JSON)
	},
)

// UsersLogin logs a user into the service
var UsersLogin = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Parse the user from the body
		user := parseUser(r)

		// Login the user
		status, message, loginToken := models.LoginUser(user)

		// Return the response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"token":   loginToken,
		})
		w.Write(JSON)
	},
)

// UsersProfile returns the current user's information
var UsersProfile = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Get user information
		status, message, retrievedUser := models.GetUser(currentUserId)

		// Return response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"user":    retrievedUser,
		})
		w.Write(JSON)
	},
)

// UsersIndex shows the information of all users
var UsersIndex = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

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
	},
)

// UsersSearch searches for a user
var UsersSearch = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get body parameters
		b, _ := ioutil.ReadAll(r.Body)
		params := make(map[string]interface{})
		json.Unmarshal(b, &params)

		// Search for users with the specified parameters
		status, message, retrievedUsers := models.SearchUsers(params)

		// Return response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"users":   retrievedUsers,
		})
		w.Write(JSON)
	},
)

// UsersShow retrieves the information of a specific user
var UsersShow = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get the specified user
		status, message, retrievedUser := models.GetUser(vars["id"])

		// Return response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"user":    retrievedUser,
		})
		w.Write(JSON)
	},
)

// UsersUpdate updates the current user
var UsersUpdate = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get body parameters
		b, _ := ioutil.ReadAll(r.Body)
		params := make(map[string]interface{})
		json.Unmarshal(b, &params)

		// Get request parameters
		vars := mux.Vars(r)

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Check if the user has permission to update the specified user
		if vars["id"] != currentUserId {
			JSON, _ := json.Marshal(map[string]interface{}{
				"status":  "error",
				"message": "You do not have permission to update this user",
				"user":    models.User{},
			})
			w.Write(JSON)
			return
		}

		// Update the user
		status, message, updatedUser := models.UpdateUser(vars["id"], params)

		// Return response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"user":    updatedUser,
		})
		w.Write(JSON)
	},
)

// UsersDelete deletes a specific user
var UsersDelete = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Check if the user has permission to delete the specified user
		if vars["id"] != currentUserId {
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
	},
)

// UsersAttendance shows the events that the current user is attending
var UsersAttendance = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get request parameters
		vars := mux.Vars(r)

		// Get the current user's ID
		currentUserId := parseCurrentUser(r)

		// Check if the user has permission to view the
		// specified user's event attendance
		if vars["id"] != currentUserId {
			JSON, _ := json.Marshal(map[string]interface{}{
				"status": "error",
				"message": "You do not have permission to" +
					" view this user's event attendance",
			})
			w.Write(JSON)
			return
		}

		// Get events the user is going to
		status, message, retrievedEvents := models.GetUserAttendance(
			vars["id"])

		// Return response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"message": message,
			"events":  retrievedEvents,
		})
		w.Write(JSON)
	},
)

// LoginWithFacebook logs in with Facebook
var LoginWithFacebook = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Get body parameters
		b, _ := ioutil.ReadAll(r.Body)
		JSONData, _ := jason.NewObjectFromBytes(b)
		token, _ := JSONData.GetString("token")

		// Get Facebook response
		response, _ := http.Get(
			"https://graph.facebook.com/me?access_token=" +
				token +
				"&fields=email,first_name,last_name,id,friends")
		defer response.Body.Close()

		// Get user details
		body, _ := ioutil.ReadAll(response.Body)
		user, _ := jason.NewObjectFromBytes([]byte(body))
		firstName, _ := user.GetString("first_name")
		lastName, _ := user.GetString("last_name")
		email, _ := user.GetString("email")
		fbId, _ := user.GetString("id")

		// Internal Facebook login
		status, message, appToken := models.ProcessFBLogin(
			firstName,
			lastName,
			email,
			fbId)

		// Return a response
		JSON, _ := json.Marshal(map[string]interface{}{
			"status":    status,
			"message":   message,
			"fb_token":  token,
			"app_token": appToken,
		})
		w.Write(JSON)
	},
)
