package main

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/config"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

func main() {

	// Initialize the database
	db.InitDB()

	// Initialize the router
	n := config.InitRouter()

	// Get the port
	port := config.GetPort()

	// Start the server
	utilities.Sugar.Infof("Started server on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), n)
	if err != nil {
		utilities.Logger.Fatal(err.Error())
	}
}
