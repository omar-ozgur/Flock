package config

import (
	"database/sql"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
)

// db holds application database information
type db struct{}

// init initializes the database
func (*db) init() {

	// Open the database
	sqlDb, err := sql.Open("postgres", utilities.DB_INFO)
	utilities.CheckErr(err)

	// Set the model db object
	models.SetDb(sqlDb)
}
