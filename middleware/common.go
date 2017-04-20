package middleware

import (
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
)

func CustomMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	utilities.Logger.Info("Started request")
	next(rw, r)
	utilities.Logger.Info("Got response")
}
