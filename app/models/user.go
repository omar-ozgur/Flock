package models

import (
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

	// Encrypt the text
	status, message, hash := EncryptText(user.Password)
	if status == "success" {
		reflections.SetField(user, "Password", string(hash))
	}

	return status, message
}

// CreateUser creates a new user
func CreateUser(user User) (status string, message string, createdUser User) {

	// Check if the user is valid
	status, message = CheckValid(user)
	if status != "success" {
		return status, message, User{}
	}

	// Encrypt password
	status, message = user.EncryptPassword()
	if status != "success" {
		return status, message, User{}
	}

	// Create the user in the database
	Db.Create(&user)
	createdUser = user
	if createdUser.ID == 0 {
		return "error", "Failed to create user", User{}
	}

	return "success", "New user created", createdUser
}

// GetUser gets a specific user
func GetUser(id string) (status string, message string, retrievedUser User) {

	// Find the user
	Db.First(&retrievedUser, id)
	if retrievedUser.ID == 0 {
		return "error", "Failed to retrieve user", User{}
	}

	return "success", "Retrieved user", retrievedUser
}

// GetUsers gets all users
func GetUsers() (status string, message string, retrievedUsers []User) {

	// Find all users
	Db.Find(&retrievedUsers)
	if len(retrievedUsers) <= 0 {
		return "error", "Failed to retrieve users", nil
	}

	return "success", "Retrieved users", retrievedUsers
}

// SearchUsers searches for users
func SearchUsers(params map[string]interface{}) (status string, message string, retrievedUsers []User) {

	// Search for users
	Db.Where(params).First(&retrievedUsers)
	if len(retrievedUsers) <= 0 {
		return "error", "Failed to retrieve users", nil
	}

	return "success", "Retrieved users", retrievedUsers
}

// UpdateUser updates a user
func UpdateUser(id string, params map[string]interface{}) (status string, message string, updatedUser User) {

	// Find the existing user
	status, message, retrievedUser := GetUser(id)
	if status != "success" {
		return status, message, User{}
	}

	// Set changed parameters
	for key, value := range params {

		// Check for password change
		if key == "Password" {
			continue
		}

		// Set updated field values
		err := reflections.SetField(&retrievedUser, key, value)
		if err != nil {
			return "error", err.Error(), User{}
		}
	}

	// Set the password if it was changed
	var changedPassword = false
	if _, exists := params["Password"]; exists {
		retrievedUser.Password = params["Password"].(string)
		changedPassword = true
	}

	// Check if the user is valid
	status, message = CheckValid(retrievedUser)
	if status != "success" {
		return status, message, User{}
	}

	// Encrypt the password if it was changed
	if changedPassword {
		status, message := retrievedUser.EncryptPassword()
		if status != "success" {
			return status, message, User{}
		}
	}

	// Update the user
	Db.Model(&retrievedUser).Updates(retrievedUser)

	// Get the updated user
	Db.First(&updatedUser, id)

	return "success", "Updated user", updatedUser
}

// DeleteUser deletes a user
func DeleteUser(id string) (status string, message string) {

	// Find the user
	status, message, retrievedUser := GetUser(id)
	if status != "success" {
		return status, message
	}

	// Find related events
	var events []Event
	Db.Model(&retrievedUser).Related(&events)

	// Delete related events
	for _, event := range events {
		Db.Unscoped().Delete(&event)
	}

	// Delete attendances
	Db.Model(&retrievedUser).Association("Attendances").Clear()

	// Delete the user
	Db.Unscoped().Delete(&retrievedUser)

	return "success", "Deleted user"
}

// CreateAttendance adds the specified user as an attendee for the event
func CreateAttendance(userId string, eventId string) (status string, message string) {

	// Get the user
	status, message, retrievedUser := GetUser(userId)
	if status != "success" {
		return status, message
	}

	// Get the event
	status, message, retrievedEvent := GetEvent(eventId)
	if status != "success" {
		return status, message
	}

	// Create attendance
	Db.Model(&retrievedUser).Association("Attendances").Append(retrievedEvent)

	return "success", "Attendance was successfully recorded"
}

// GetUserAttendance gets events that a specific user is attending
func GetUserAttendance(userId string) (status string, message string, retrievedEvents []Event) {

	// Get the user
	status, message, retrievedUser := GetUser(userId)
	if status != "success" {
		return status, message, nil
	}

	// Find events
	Db.Model(&retrievedUser).Association("Attendances").Find(&retrievedEvents)

	return "success", "Retrieved events", retrievedEvents
}

// DeleteAttendance removes the specified user from the event's attendee list
func DeleteAttendance(userId string, eventId string) (status string, message string) {

	// Get the user
	status, message, retrievedUser := GetUser(userId)
	if status != "success" {
		return status, message
	}

	// Get the event
	status, message, retrievedEvent := GetEvent(eventId)
	if status != "success" {
		return status, message
	}

	// Delete attendance
	Db.Model(&retrievedUser).Association("Attendances").Delete(retrievedEvent)

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
	Db.Where("email = ?", user.Email).First(&retrievedUser)

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password))
	if err != nil {
		return "error", "Error while checking password", ""
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = retrievedUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(utilities.FLOCK_TOKEN_SECRET)

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
