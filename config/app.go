package config

import (
	"github.com/omar-ozgur/flock-api/app/models"
)

// Application maintains pipelined application logic
type Application struct {
	db     *db
	router *router
	server *server
}

// NewApplication initializes a new application
func NewApplication() *Application {
	application := Application{}

	// Create a new database
	application.db = NewDb()

	// Create a new router
	application.router = NewRouter()

	// Create a new server
	application.server = NewServer()

	// Initialize models
	models.Init()

	return &application
}

// Start starts the application
func (application *Application) Start() {

	// Start the server
	application.server.start(application.router.middleware.Negroni)
}
