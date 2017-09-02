package utilities

import (
	"os"
)

// DEFAULT_PORT is the default port used in the application
var DEFAULT_PORT = "3000"

// PORT is the port that is actually being used by the application
var PORT = func() string {

	// Attempt to retrieve the port from the environment
	port := os.Getenv("PORT")

	// Use the default port if the environment variable was blank
	if port == "" {
		port = DEFAULT_PORT
		Sugar.Infof(
			"No PORT environment variable detected, defaulting to port %s",
			port,
		)
	}

	return port
}
