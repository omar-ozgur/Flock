package middleware

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Middleware holds application middleware information
type Middleware struct {
	Negroni *negroni.Negroni
}

// NewMiddleware initializes application middleware
func NewMiddleware(r *mux.Router) *Middleware {
	middleware := Middleware{}

	// Create new negroni middleware
	middleware.Negroni = negroni.New(
		negroni.HandlerFunc(loggingMiddleware),
		negroni.NewLogger(),
	)
	middleware.Negroni.UseHandler(r)

	return &middleware
}
