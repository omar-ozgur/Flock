package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.PagesIndex)

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
