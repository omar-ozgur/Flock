package controllers

import (
	"net/http"
)

func PagesIndex(w http.ResponseWriter, r *http.Request) {
	a := new(PagesAttributes)
	a.Username = "Omar"
	templates.ExecuteTemplate(w, "index.html", a)
}

type PagesAttributes struct {
	Username string
}
