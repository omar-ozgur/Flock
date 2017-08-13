package models

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oleiade/reflections.v1"
	"reflect"
	"strconv"
	"time"
)

// User model
type User struct {
	Id           int       `valid:"-"`
	First_name   string    `valid:"required"`
	Last_name    string    `valid:"required"`
	Email        string    `valid:"email,required"`
	Fb_id        string    `valid:"-"`
	Password     []byte    `valid:"required"`
	Time_created time.Time `valid:"-"`
}

// Parameters that are created automatically
var userAutoParams = map[string]bool{"Id": true, "Time_created": true}

// Parameters that must be unique
var userUniqueParams = map[string]bool{"Email": true, "Fb_id": true}

// Parameters that are required
var userRequiredParams = map[string]bool{"First_name": true, "Last_name": true, "Email": true, "Password": true}

// Create a user
func CreateUser(user User) (status string, message string, createdUser User) {

	// Encrypt password
	hash, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return "error", fmt.Sprintf("Failed to encrypt password: %s", err.Error()), User{}
	}
	reflections.SetField(&user, "Password", hash)

	// Get user fields
	fields := reflect.ValueOf(user)

	// Validate user
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate user: %s", err.Error()), User{}
	}

	// Create a map of unique parameter values
	uniqueMap := make(map[string]interface{})
	for key, _ := range userUniqueParams {

		// Add parameter value to unique map
		fieldValue, err := reflections.GetField(&user, key)
		if err != nil {
			return "error", fmt.Sprintf("Failed to get field: %s", err.Error()), User{}
		}
		uniqueMap[key] = fieldValue
	}

	// Search for users with the same unique parameters
	status, message, retrievedUsers := SearchUsers(uniqueMap, "AND")
	if status != "success" {
		return "error", "Failed to check user uniqueness", User{}
	} else if retrievedUsers != nil {
		return "error", "User is not unique", User{}
	}

	// Set present column names
	var fieldsStr, valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < fields.NumField(); i++ {

		// Get field name and value
		fieldName := fields.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", fields.Field(i).Interface())

		// Skip the field if it is automatically set by the database
		if userAutoParams[fieldName] {
			continue
		}

		// Add a comma if it is not the first field
		if !first {
			fieldsStr.WriteString(", ")
			valuesStr.WriteString(", ")
		} else {
			first = false
		}

		// Add field name and value to lists
		fieldsStr.WriteString(fieldName)
		valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
		values = append(values, fieldValue)
		parameterIndex += 1
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", utilities.USERS_TABLE))

	// Finish query string
	queryStr.WriteString(fmt.Sprintf("%s) VALUES(%s) RETURNING id;", fieldsStr.String(), valuesStr.String()))

	// Log query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), User{}
	}

	// Execute query
	err = stmt.QueryRow(values...).Scan(&user.Id)
	if err != nil {
		return "error", fmt.Sprintf("Failed to create new user: %s", err.Error()), User{}
	}

	// Get created user
	status, _, createdUser = GetUser(fmt.Sprintf("%v", user.Id))
	if status != "success" {
		return "error", "Failed to retrieve created user", User{}
	}

	return "success", "New user created", createdUser
}

// Login a user
func LoginUser(user User) (status string, message string, createdToken string) {

	// Check login parameter presence
	if user.Email == "" {
		return "error", "Email cannot be blank", ""
	} else if len(user.Password) == 0 {
		return "error", "Password cannot be blank", ""
	}

	// Create a query to find a user
	var foundUser User
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE email=$1;", utilities.USERS_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", user.Email)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), ""
	}

	// Execute the query
	row := stmt.QueryRow(user.Email)

	// Parse user values
	err = row.Scan(&foundUser.Id, &foundUser.First_name, &foundUser.Last_name, &foundUser.Email, &foundUser.Fb_id, &foundUser.Password, &foundUser.Time_created)
	if err != nil {
		return "error", "Error while retrieving user", ""
	}

	// Check password
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "error", "Error while encrypting password", ""
	}
	err = bcrypt.CompareHashAndPassword(hash, user.Password)
	if err != nil {
		return "error", "Error while checking password", ""
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = foundUser.Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(utilities.FLOCK_TOKEN_SECRET)

	return "success", "Login token generated", tokenString
}

// Get a user
func GetUser(id string) (status string, message string, retrievedUser User) {

	// Create a query to get the user
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", utilities.USERS_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), User{}
	}

	// Execute the query
	row := stmt.QueryRow(id)

	// Get user info
	var user User
	err = row.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
	if err != nil {
		return "error", "Failed to retrieve user information", User{}
	}

	return "success", "Retrieved user", user
}

// Get all users
func GetUsers() (status string, message string, retrievedUsers []User) {

	// Create the query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", utilities.USERS_TABLE)

	// Log the query
	utilities.Sugar.Infof("SQL Query: %s", queryStr)

	// Execute the query
	rows, err := db.DB.Query(queryStr)
	if err != nil {
		return "error", "Failed to query users", nil
	}

	// Get user info
	var users []User
	for rows.Next() {

		// Parse a user
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
		if err != nil {
			return "error", "Failed to retrieve user information", nil
		}

		// Add the user to the users array
		users = append(users, user)
	}

	return "success", "Retrieved users", users
}

// Update user
func UpdateUser(id string, user User) (status string, message string, updatedUser User) {

	// Get user fields
	fields := reflect.ValueOf(user)
	if fields.NumField() <= 0 {
		return "error", "Invalid number of fields", user
	}

	// Create a map of unique parameter values
	uniqueMap := make(map[string]interface{})
	for key, _ := range userUniqueParams {

		// Get the field value
		fieldValue, err := reflections.GetField(&user, key)

		// Skip the field if there was an error, or if it is empty
		if err != nil || reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}

		// Add the field value to the unique map
		uniqueMap[key] = fieldValue
	}

	// Search for users that already have the unique values
	if len(uniqueMap) > 0 {
		status, _, retrievedUsers := SearchUsers(uniqueMap, "OR")
		if status != "success" {
			return "error", "Failed to check user uniqueness", User{}
		} else if retrievedUsers != nil {
			return "error", "User is not unique", User{}
		}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", utilities.USERS_TABLE))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < fields.NumField(); i++ {

		// Get field name and value
		fieldName := fields.Type().Field(i).Name
		fieldValue := fields.Field(i).Interface()

		// Skip the field if it is created automatically by the database
		if userAutoParams[fieldName] {
			continue
		}

		// Skip the parameter if it is empty
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}

		// Add delimiters
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}

		// Check if the field is a password
		if fieldName == "Password" {

			// Encrypt the password and add it to the query
			hash, err := bcrypt.GenerateFromPassword([]byte(fieldValue.(string)), bcrypt.DefaultCost)
			if err != nil {
				return "error", "Failed to encrypt password", User{}
			}
			queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
			values = append(values, hash)
		} else {

			// Add the field value to the query
			queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
			values = append(values, fieldValue)
		}

		parameterIndex += 1
	}

	// Add the id to the values
	values = append(values, id)

	// Finish query
	queryStr.WriteString(fmt.Sprintf(" WHERE id=$%d;", parameterIndex))

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), User{}
	}

	// Execute the query
	_, err = stmt.Exec(values...)
	if err != nil {
		return "error", fmt.Sprintf("Failed to update user: %s", err.Error()), User{}
	}

	// Get the updated user
	status, message, retrievedUser := GetUser(id)
	if status == "success" {
		return "success", "Updated user", retrievedUser
	} else {
		return "error", "Failed to retrieve updated user", User{}
	}
}

// Delete a user
func DeleteUser(id string) (status string, message string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", utilities.USERS_TABLE)

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}

	// Execute the query
	_, err = stmt.Exec(id)
	if err != nil {
		return "error", "Failed to delete user"
	}

	// Delete the user's events
	i, _ := strconv.Atoi(id)
	status, _, retrievedEvents := SearchEvents(Event{User_id: i})
	if status == "success" {
		for _, event := range retrievedEvents {
			DeleteEvent(fmt.Sprintf("%d", event.Id))
		}
	}

	// Delete the user's attendance
	attendeeParams := make(map[string]interface{})
	attendeeParams["User_id"] = i
	status, _, retrievedAttendees := SearchAttendees(attendeeParams, "AND")
	if status == "success" {
		for _, attendee := range retrievedAttendees {
			DeleteAttendee(attendee)
		}
	}

	return "success", "Deleted user"
}

// Search for users
func SearchUsers(parameters map[string]interface{}, operator string) (status string, message string, retrievedUsers []User) {

	// Create query string
	var queryStr bytes.Buffer

	// Create and execute query
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", utilities.USERS_TABLE))

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for key, value := range parameters {

		// Add delimiters
		if !first {
			queryStr.WriteString(fmt.Sprintf(" %s ", operator))
		} else {
			queryStr.WriteString(" ")
			first = false
		}

		// Add the field to the query
		queryStr.WriteString(fmt.Sprintf("%v=$%d", key, parameterIndex))
		values = append(values, value)

		parameterIndex += 1
	}

	// Finish the query
	queryStr.WriteString(";")

	// Log the query and values
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)

	// Prepare the query
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}

	// Execute the query
	rows, err := stmt.Query(values...)
	if err != nil {
		return "error", "Failed to query users", nil
	}

	// Parse users
	var users []User
	for rows.Next() {

		// Parse user
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
		if err != nil {
			return "error", "Failed to retrieve user information", nil
		}

		// Add user to the list
		users = append(users, user)
	}

	return "success", "Retrieved users", users
}

func GetUserAttendance(id string) (status string, message string, retrievedEvents []Event) {

	// Find attendee objects
	attendeeParams := make(map[string]interface{})
	attendeeParams["User_id"] = id
	status, message, retrievedAttendees := SearchAttendees(attendeeParams, "AND")
	if status != "success" {
		return "error", "Failed to check search attendees", nil
	}

	// Find events associated with attendee objects
	var events []Event
	for i := range retrievedAttendees {

		// Search for events based on the attendee object
		event := Event{Id: retrievedAttendees[i].Event_id}
		status, _, retrievedEvents := SearchEvents(event)
		if status != "success" {
			return "error", "Failed to retrieve events based on attendee information", nil
		}

		// Add the found events to the events list
		events = append(events, retrievedEvents...)
	}

	return "success", "Retrieved events", events
}

// Process Facebook login
func ProcessFBLogin(first_name string, last_name string, email string, fb_id string) (status string, message string, accessToken string) {

	// Create a search query
	searchQuery := make(map[string]interface{})
	searchQuery["first_name"] = first_name
	searchQuery["last_name"] = last_name
	searchQuery["email"] = email
	searchQuery["fb_id"] = fb_id

	// Search for users with the specified fields
	status, message, retrievedUsers := SearchUsers(searchQuery, "AND")
	if status != "success" {
		return "error", message, ""
	}

	// If the user doesn't exist, create it
	if len(retrievedUsers) == 0 {

		// Create a user object
		var user User
		user.First_name = searchQuery["first_name"].(string)
		user.Last_name = searchQuery["last_name"].(string)
		user.Email = searchQuery["email"].(string)
		user.Fb_id = searchQuery["fb_id"].(string)
		user.Password = []byte("Facebook_User")

		// Create the user
		status, message, createdUser := CreateUser(user)
		if status != "success" {
			return "error", message, ""
		}

		// Add the new user to the list of users
		retrievedUsers = append(retrievedUsers, createdUser)
	}

	// Login the user
	return LoginUser(retrievedUsers[0])
}
