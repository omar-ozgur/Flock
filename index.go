package main

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/config"
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

func main() {

	n := config.InitRouter()

	utilities.Sugar.Infof("Started server on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), n)
	if err != nil {
		utilities.Logger.Fatal(err.Error())
	}
}
