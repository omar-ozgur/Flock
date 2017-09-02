package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/utilities"
	"github.com/urfave/negroni"
	"net/http"
)

// initMiddleware initializes application middleware
func Init(r *mux.Router) *negroni.Negroni {

	// Create new negroni middleware
	negroni := negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.NewLogger(),
	)
	negroni.UseHandler(r)

	return negroni
}

// Logging middleware
func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Indicate the start of a request
	utilities.Logger.Info("Started request")

	next(rw, r)

	// Indicate the end of a request
	utilities.Logger.Info("Got response")

	// Print a newline to separate request logs
	fmt.Println()
}
