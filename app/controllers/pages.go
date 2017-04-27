package controllers

import (
	"net/http"
)

func PagesIndex(w http.ResponseWriter, r *http.Request) {
	page := new(PagesAttributes)
	page.Name = "home"
	templates.ExecuteTemplate(w, "index.html", page)
}

type PagesAttributes struct {
	Name string
}
