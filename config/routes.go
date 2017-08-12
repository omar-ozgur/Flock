package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

// Initialize the router
func InitRouter() (n *negroni.Negroni) {

	// Create an authorization handler
	authorizationHandler := middleware.JWTMiddleware.Handler

	// Create a new router
	router := mux.NewRouter()

	// Create user routes
	router.Handle("/signup", controllers.UsersCreate).Methods("POST")
	router.Handle("/login", controllers.UsersLogin).Methods("POST")
	router.Handle("/profile", authorizationHandler(controllers.UsersProfile)).Methods("Get")
	router.Handle("/users", controllers.UsersIndex).Methods("GET")
	router.Handle("/users/search", controllers.UsersSearch).Methods("POST")
	router.Handle("/users/{id}", controllers.UsersShow).Methods("GET")
	router.Handle("/users/{id}", authorizationHandler(controllers.UsersUpdate)).Methods("PUT")
	router.Handle("/users/{id}", authorizationHandler(controllers.UsersDelete)).Methods("DELETE")
	router.Handle("/users/{id}/attendance", authorizationHandler(controllers.UsersAttendance)).Methods("GET")

	// Create event routes
	router.Handle("/events", controllers.EventsIndex).Methods("GET")
	router.Handle("/events", authorizationHandler(controllers.EventsCreate)).Methods("POST")
	router.Handle("/events/search", controllers.EventsSearch).Methods("POST")
	router.Handle("/events/{id}", controllers.EventsShow).Methods("GET")
	router.Handle("/events/{id}", authorizationHandler(controllers.EventsUpdate)).Methods("PUT")
	router.Handle("/events/{id}", authorizationHandler(controllers.EventsDelete)).Methods("DELETE")
	router.Handle("/events/{id}/attendees", controllers.EventsAttendees).Methods("GET")
	router.Handle("/events/{id}/attend", authorizationHandler(controllers.EventsAttend)).Methods("POST")
	router.Handle("/events/{id}/attendance", authorizationHandler(controllers.EventsDeleteAttendance)).Methods("DELETE")

	// Create login routes
	router.HandleFunc("/loginWithFacebook", controllers.LoginWithFacebook).Methods("POST")

	// Integrate middleware
	n = negroni.New(negroni.HandlerFunc(middleware.LoggingMiddleware), negroni.NewLogger())
	n.UseHandler(router)

	return
}
