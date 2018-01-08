package config

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/utilities"
	"github.com/urfave/negroni"
	"net/http"
)

// server holds application server information
type server struct{}

// NewServer initialized a new server
func NewServer() *server {
	server := server{}
	return &server
}

// start starts the application server
func (*server) start(negroni *negroni.Negroni) {
	utilities.Sugar.Infof("Started server on port %s\n", utilities.PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", utilities.PORT), negroni)
	if err != nil {
		utilities.Logger.Fatal(err.Error())
	}
}
