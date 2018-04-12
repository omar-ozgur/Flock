package models

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/omar-ozgur/flock-api/utilities"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oleiade/reflections.v1"
	"time"
)

// User holds data for application users
type User struct {
	gorm.Model
	First_name string `gorm:"not null" valid:"required"`
	Last_name  string `gorm:"not null" valid:"required"`
	Email      string `gorm:"unique;not null" valid:"email,required"`
	Fb_id      string `gorm:"-"`
	Password   string `gorm:"not null" valid:"required"`
}

// InitUsers initializes users
func InitUsers() {
	Db.AutoMigrate(&User{})
}

// encryptPassword encrypts a password
func encryptPassword(password string) (status string, message string, hash []byte) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "error", fmt.Sprintf("Failed to encrypt password: %s", err.Error()), nil
	}

	return "success", "", hash
}

// encryptPassword encrypts a user's password
func (user *User) encryptPassword() (status string, message string) {

	status, message, hash := encryptPassword(user.Password)
	if status == "success" {
		reflections.SetField(user, "Password", string(hash))
	}

	return status, message
}

// checkValid checks if a user is valid
func (user *User) checkValid() (status string, message string) {

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate user: %s", err.Error())
	}

	return "success", ""
}

// CreateUser creates a new user
func CreateUser(user User) (status string, message string, createdUser User) {

	// Check if the user is valid
	status, message = user.checkValid()
	if status != "success" {
		return status, message, User{}
	}

	// Encrypt password
	status, message = user.encryptPassword()
	if status != "success" {
		return status, message, User{}
	}

	// Create the user in the database
	Db.Create(&user)

	// Get the created user
	Db.First(&createdUser, user.ID)
	if createdUser.ID == 0 {
		return "error", "Failed to retrieve created user", User{}
	}

	return "success", "New user created", createdUser
}

// GetUser gets a specific user
func GetUser(id string) (status string, message string, retrievedUser User) {

	// Find the user
	Db.First(&retrievedUser, id)
	if retrievedUser.ID == 0 {
		return "error", "Could not find user", User{}
	}

	return "success", "Retrieved user", retrievedUser
}

// GetUsers gets all users
func GetUsers() (status string, message string, retrievedUsers []User) {

	// Find all users
	Db.Find(&retrievedUsers)
	if len(retrievedUsers) <= 0 {
		return "error", "Could not find any users", nil
	}

	return "success", "Retrieved users", retrievedUsers
}

// UpdateUser updates a user
func UpdateUser(id string, params map[string]interface{}) (status string, message string, updatedUser User) {

	// Find the existing user
	status, message, foundUser := GetUser(id)
	if status != "success" {
		return status, message, User{}
	}

	// Set changed parameters
	for key, value := range params {
		if key == "Password" {
			continue
		}

		err := reflections.SetField(&foundUser, key, value)
		if err != nil {
			return "error", err.Error(), User{}
		}
	}

	// Set the password if it was changed
	var changedPassword = false
	if _, ok := params["Password"]; ok {
		foundUser.Password = params["Password"].(string)
		changedPassword = true
	}

	// Check if the user is valid
	status, message = foundUser.checkValid()
	if status != "success" {
		return status, message, User{}
	}

	// Encrypt the password if it was changed
	if changedPassword {
		status, message := foundUser.encryptPassword()
		if status != "success" {
			return status, message, User{}
		}
	}

	// Update the user
	err := Db.Model(&foundUser).Updates(foundUser).Error
	if err != nil {
		return "error", "Failed to update user", User{}
	}

	// Get the updated user
	Db.First(&updatedUser, id)
	if updatedUser.ID == 0 {
		return "error", "Failed to retrieve updated user", User{}
	}

	return "success", "Updated user", updatedUser
}

// DeleteUser deletes a user
func DeleteUser(id string) (status string, message string) {

	// Find the user
	var user User
	Db.First(&user, id)

	// Delete the user
	Db.Unscoped().Delete(&user)

	return "success", "Deleted user"
}

// LoginUser logs in a user
func LoginUser(user User) (status string, message string, createdToken string) {

	// Check login parameter presence
	if user.Email == "" || user.Password == "" {
		return "error", "Email and password cannot be blank", ""
	}

	// Find the user
	var foundUser User
	Db.Where("email = ?", user.Email).First(&foundUser)

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return "error", "Error while checking password", ""
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = foundUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(utilities.FLOCK_TOKEN_SECRET)

	return "success", "Login token generated", tokenString
}

// SearchUsers searches for users
func SearchUsers(parameters map[string]interface{}) (status string, message string, retrievedUsers []User) {

	Db.Where(parameters).First(&retrievedUsers)

	return "success", "Retrieved users", retrievedUsers
}

// GetUserAttendance gets events that a specific user is attending
func GetUserAttendance(id string) (status string, message string, retrievedEvents []Event) {

	// Find attendee objects
	attendeeParams := make(map[string]interface{})
	attendeeParams["User_id"] = id
	status, message, retrievedAttendees := SearchAttendees(attendeeParams)
	if status != "success" {
		return "error", "Failed to check search attendees", nil
	}

	// Find events associated with attendee objects
	var events []Event
	for i := range retrievedAttendees {

		// Search for events based on the attendee object
		eventParams := make(map[string]interface{})
		eventParams["ID"] = retrievedAttendees[i].Event_id
		status, _, retrievedEvents := SearchEvents(eventParams)
		if status != "success" {
			return "error", "Failed to retrieve events based on attendee information", nil
		}

		// Add the found events to the events list
		events = append(events, retrievedEvents...)
	}

	return "success", "Retrieved events", events
}

// ProcessFBLogin processes Facebook login
func ProcessFBLogin(first_name string, last_name string, email string, fb_id string) (status string, message string, accessToken string) {

	// Create a search query
	searchQuery := make(map[string]interface{})
	searchQuery["first_name"] = first_name
	searchQuery["last_name"] = last_name
	searchQuery["email"] = email
	searchQuery["fb_id"] = fb_id

	// Search for users with the specified fields
	status, message, retrievedUsers := SearchUsers(searchQuery)
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
		user.Password = "Facebook_User"

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
