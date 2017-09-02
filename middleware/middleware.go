package middleware

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Middleware holds application middleware information
type Middleware struct {
	Negroni *negroni.Negroni
}

// init initializes application middleware
func (m *Middleware) Init(r *mux.Router) {

	// Create new negroni middleware
	m.Negroni = negroni.New(
		negroni.HandlerFunc(loggingMiddleware),
		negroni.NewLogger(),
	)
	m.Negroni.UseHandler(r)
}
