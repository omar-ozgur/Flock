package controllers

import (
	"github.com/omar-ozgur/flock-api/utilities"
	"net/http"
	//"fmt"
)

func PagesIndex(w http.ResponseWriter, r *http.Request) {
	page := new(PagesAttributes)
	page.Name = "home"

	url := utilities.FbConfig.AuthCodeURL("")

	//fmt.Println(url);

	page.URL = url
	templates.ExecuteTemplate(w, "index.html", page)
}

type PagesAttributes struct {
	Name string
	URL  string
}
