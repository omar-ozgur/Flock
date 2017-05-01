package main

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/config"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

func main() {

	n := config.InitRouter()

	db.InitDB()

	port := config.GetPort()
	utilities.Sugar.Infof("Started server on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), n)
	if err != nil {
		utilities.Logger.Fatal(err.Error())
	}
}
