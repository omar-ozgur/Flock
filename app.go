package main

import (
	"github.com/omar-ozgur/flock-api/config"
)

// main is the entry-point for the flock API
func main() {

	// Create a new application
	application := config.NewApplication()

	// Start the application
	application.Start()
}
