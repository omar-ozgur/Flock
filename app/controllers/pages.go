package controllers

import (
	"net/http"
	"github.com/omar-ozgur/flock-api/utilities"
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
	URL string
}
