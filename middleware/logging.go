package middleware

import (
	"fmt"
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

// loggingMiddleware creates standard logging middleware
// Logs the start and end of a request
func loggingMiddleware(
	rw http.ResponseWriter,
	r *http.Request,
	next http.HandlerFunc) {

	// Indicate the start of a request
	utilities.Logger.Info("Started request")

	next(rw, r)

	// Indicate the end of a request
	utilities.Logger.Info("Got response")

	// Print a newline to separate request logs
	fmt.Println()
}
