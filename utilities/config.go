package utilities

import (
	"os"
)

// General server configuration values
var DEFAULT_PORT = "3000"
var PORT = getPort()

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
		Sugar.Infof(
			"No PORT environment variable detected, defaulting to port %s",
			port,
		)
	}
	return port
}
