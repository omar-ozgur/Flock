package models

import (
	"github.com/jinzhu/gorm"
	"github.com/omar-ozgur/flock-api/utilities"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oleiade/reflections.v1"
	"time"
)

// User holds data for application users
type User struct {
	gorm.Model
	First_name  string   `gorm:"not null" valid:"required"`
	Last_name   string   `gorm:"not null" valid:"required"`
	Email       string   `gorm:"unique;not null" valid:"email,required"`
	Fb_id       string   `gorm:"-"`
	Password    string   `gorm:"not null" valid:"required"`
	Attendances []*Event `gorm:"many2many:attendees;"`
	Events      []Event
}

// InitUsers initializes users
func InitUsers() {
	Db.AutoMigrate(&User{})
}

// EncryptPassword encrypts a user's password
func (user *User) EncryptPassword() (status string, message string) {
	var hash []byte
	if status, message, hash = EncryptText(user.Password); status != "success" {
		return status, message
	}

	reflections.SetField(user, "Password", string(hash))
	return "success", "Successfully encrypted the password"
}

// CreateUser creates a new user
func CreateUser(user User) (status string, message string, createdUser User) {

	// Begin the transaction
	tx := Db.Begin()

	// Check if the user is valid
	if status, message = CheckValid(user); status != "success" {
		return status, message, User{}
	}

	// Encrypt password
	if status, message = user.EncryptPassword(); status != "success" {
		return status, message, User{}
	}

	// Create the user in the database
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to create user", User{}
	}

	// Commit the transaction
	tx.Commit()
	return "success", "New user created", user
}

// GetUser gets a specific user
func GetUser(id string) (status string, message string, retrievedUser User) {

	// Find the user
	if err := Db.First(&retrievedUser, id).Error; err != nil {
		return "error", "Failed to retrieve user", User{}
	}

	return "success", "Retrieved user", retrievedUser
}

// GetUsers gets all users
func GetUsers() (status string, message string, retrievedUsers []User) {

	// Find all users
	if err := Db.Find(&retrievedUsers).Error; err != nil {
		return "error", "Failed to retrieve users", nil
	}

	return "success", "Retrieved users", retrievedUsers
}

// SearchUsers searches for users
func SearchUsers(params map[string]interface{}) (status string, message string, retrievedUsers []User) {

	// Search for users
	if err := Db.Where(params).First(&retrievedUsers).Error; err != nil {
		return "error", "Failed to retrieve users", nil
	}

	return "success", "Retrieved users", retrievedUsers
}

// UpdateUser updates a user
func UpdateUser(id string, params map[string]interface{}) (status string, message string, updatedUser User) {

	// Begin the transaction
	tx := Db.Begin()

	// Find the existing user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(id); status != "success" {
		return status, message, User{}
	}

	// Set changed parameters
	for key, value := range params {
		if err := reflections.SetField(&retrievedUser, key, value); err != nil {
			return "error", "Failed to set updated field", User{}
		}
	}

	// Check if the user is valid
	if status, message = CheckValid(retrievedUser); status != "success" {
		return status, message, User{}
	}

	// Encrypt the password if it was changed
	if _, exists := params["Password"]; exists {
		if status, message := retrievedUser.EncryptPassword(); status != "success" {
			return status, message, User{}
		}
	}

	// Update the user
	if err := tx.Model(&retrievedUser).Updates(retrievedUser).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to update user", User{}
	}

	// Get the updated user
	if err := tx.First(&updatedUser, id).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to retrieve updated user", User{}
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Updated user", updatedUser
}

// DeleteUser deletes a user
func DeleteUser(id string) (status string, message string) {

	// Begin the transaction
	tx := Db.Begin()

	// Find the user
	var retrievedUser User
	if status, message, retrievedUser = GetUser(id); status != "success" {
		return status, message
	}

	// Find related events
	var events []Event
	if err := tx.Model(&retrievedUser).Related(&events).Error; err != nil {
		return "error", "Failed to retrieve related events"
	}

	// Delete related events
	for _, event := range events {
		if err := tx.Unscoped().Delete(&event).Error; err != nil {
			tx.Rollback()
			return "error", "Failed to delete related events"
		}
	}

	// Delete attendances
	if err := tx.Model(&retrievedUser).Association("Attendances").Clear().Error; err != nil {
		tx.Rollback()
		return "error", "Failed to delete attendances"
	}

	// Delete the user
	if err := tx.Unscoped().Delete(&retrievedUser).Error; err != nil {
		tx.Rollback()
		return "error", "Failed to delete the user"
	}

	// Commit the transaction
	tx.Commit()
	return "success", "Deleted user"
}

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

// LoginUser logs in a user
func LoginUser(user User) (status string, message string, createdToken string) {

	// Check login parameter presence
	if user.Email == "" || user.Password == "" {
		return "error", "Email and password cannot be blank", ""
	}

	// Find the user
	var retrievedUser User
	if err := Db.Where("email = ?", user.Email).First(&retrievedUser).Error; err != nil {
		return "error", "Failed to retrieve the user", ""
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password)); err != nil {
		return "error", "Error while checking password", ""
	}

	// Create a new token
	claims := make(map[string]interface{})
	claims["user_id"] = retrievedUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString := utilities.CreateToken(claims)

	return "success", "Login token generated", tokenString
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
	var retrievedUsers []User
	if status, message, retrievedUsers = SearchUsers(searchQuery); status != "success" {
		return status, message, ""
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
		var createdUser User
		if status, message, createdUser = CreateUser(user); status != "success" {
			return status, message, ""
		}

		// Add the new user to the list of users
		retrievedUsers = append(retrievedUsers, createdUser)
	}

	// Login the user
	return LoginUser(retrievedUsers[0])
}
