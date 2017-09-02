package config

import (
	"github.com/omar-ozgur/flock-api/app/models"
)

// ApplicationFacade maintains pipelined application logic
type ApplicationFacade struct{}

// Init initializes the application
func (ApplicationFacade) Init() {

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
