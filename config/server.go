package config

import (
	"github.com/omar-ozgur/flock-api/utilities"
	"os"
)

var defaultPort = 3000

func GetPort() string {
	var port = os.Getenv("PORT")

	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000"
		utilities.Sugar.Infof("No PORT environment variable detected, defaulting to port %s", port)
	}
	return port
}
