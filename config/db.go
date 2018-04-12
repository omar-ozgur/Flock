package config

import (
	"github.com/jinzhu/gorm"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
)

// db holds application database information
type db struct {
	gormDb *gorm.DB
}

// NewDb initializes a new database
func NewDb() *db {
	db := db{}

	// Open the database
	gormDb, err := gorm.Open("postgres", utilities.DB_INFO)
	utilities.CheckErr(err)
	db.gormDb = gormDb

	// Set the model db object
	models.SetDb(db.gormDb)

	return &db
}
