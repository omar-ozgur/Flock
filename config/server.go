package config

import (
	"github.com/omar-ozgur/flock-api/utilities"
	"os"
)

// Get the server port
func GetPort() string {

	// Set a default port if there is nothing in the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = utilities.DEFAULT_PORT
		utilities.Sugar.Infof("No PORT environment variable detected, defaulting to port %s", port)
	}
	return port
}
