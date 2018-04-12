package models

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Db is the database that the models will use
var Db *gorm.DB

// SetDb sets the database for the models to use
func SetDb(db *gorm.DB) {
	Db = db
}

// Init initializes all models
func Init() {
	InitUsers()
	InitEvents()
}

// EncryptText encrypts text
func EncryptText(text string) (status string, message string, hash []byte) {
	var err error
	if hash, err = bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost); err != nil {
		return "error", fmt.Sprintf("Failed to encrypt text: %s", err.Error()), nil
	}

	return "success", "Successfully encrypted the text", hash
}

// CheckValid checks if the model is valid
func CheckValid(object interface{}) (status string, message string) {
	if _, err := govalidator.ValidateStruct(object); err != nil {
		return "error", fmt.Sprintf("Failed to validate: %s", err.Error())
	}

	return "success", ""
}
