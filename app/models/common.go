package models

import (
	"github.com/jinzhu/gorm"
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
	InitAttendees()
}
