package main

import (
	"github.com/omar-ozgur/flock-api/config"
)

// main is the entry-point for the flock API
func main() {

	// Initialize the application
	application := config.Application{}
	application.Init()
}
