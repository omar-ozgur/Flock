package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.PagesIndex).Methods("GET")
	r.HandleFunc("/users", controllers.UsersIndex).Methods("GET")
	r.HandleFunc("/users", controllers.UsersCreate).Methods("POST")
	r.HandleFunc("/users/{id}", controllers.UsersShow).Methods("GET")
	r.HandleFunc("/users/{id}", controllers.UsersUpdate).Methods("PUT")
	r.HandleFunc("/users/{id}", controllers.UsersDelete).Methods("DELETE")

	r.HandleFunc("/posts", controllers.PostsIndex).Methods("GET")
	r.HandleFunc("/posts", controllers.PostsCreate).Methods("POST")
	r.HandleFunc("/posts/{id}", controllers.PostsShow).Methods("GET")
	r.HandleFunc("/posts/{id}", controllers.PostsUpdate).Methods("PUT")
	r.HandleFunc("/posts/{id}", controllers.PostsDelete).Methods("DELETE")

	r.HandleFunc("/attendees", controllers.AttendeesIndex).Methods("GET")
	r.HandleFunc("/attendees", controllers.AttendeesCreate).Methods("POST")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesShow).Methods("GET")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesUpdate).Methods("PUT")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesDelete).Methods("DELETE")

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
