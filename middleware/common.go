package middleware

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

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
