package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	authorizationHandler := middleware.JWTMiddleware.Handler

	r := mux.NewRouter()

	r.HandleFunc("/", controllers.PagesIndex).Methods("GET")

	r.Handle("/signup", controllers.UsersCreate).Methods("POST")
	r.Handle("/login", controllers.UsersLogin).Methods("POST")
	r.Handle("/profile", authorizationHandler(controllers.UsersProfile)).Methods("Get")
	r.Handle("/users", controllers.UsersIndex).Methods("GET")
	r.Handle("/users/search", controllers.UsersSearch).Methods("POST")
	r.Handle("/users/{id}", controllers.UsersShow).Methods("GET")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersUpdate)).Methods("PUT")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersDelete)).Methods("DELETE")
	r.Handle("/users/{id}/attendance", authorizationHandler(controllers.UsersAttendance)).Methods("GET")

	r.Handle("/events", controllers.EventsIndex).Methods("GET")
	r.Handle("/events", authorizationHandler(controllers.EventsCreate)).Methods("POST")
	r.Handle("/events/search", controllers.EventsSearch).Methods("POST")
	r.Handle("/events/{id}", controllers.EventsShow).Methods("GET")
	r.Handle("/events/{id}", authorizationHandler(controllers.EventsUpdate)).Methods("PUT")
	r.Handle("/events/{id}", authorizationHandler(controllers.EventsDelete)).Methods("DELETE")
	r.Handle("/events/{id}/attendees", controllers.EventsAttendees).Methods("GET")
	r.Handle("/events/{id}/attend", authorizationHandler(controllers.EventsAttend)).Methods("POST")
	r.Handle("/events/{id}/attendance", authorizationHandler(controllers.EventsDeleteAttendance)).Methods("DELETE")

	//login handling for Facebook
	r.HandleFunc("/loginWithFB", controllers.LoginWithFacebook).Methods("POST")

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
