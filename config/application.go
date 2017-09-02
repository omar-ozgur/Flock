package config

import (
	"github.com/omar-ozgur/flock-api/app/models"
)

// Application maintains pipelined application logic
type Application struct{}

// Init initializes the application
func (Application) Init() {

	// Initialize the database
	db := db{}
	db.init()

	// Initialize the router
	router := router{}
	negroni := router.init()

	// Initialize models
	models.Init()

	// Start the server
	server := server{}
	server.start(negroni)
}
