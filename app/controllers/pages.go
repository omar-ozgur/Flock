package controllers

import (
	"encoding/json"
	"github.com/omar-ozgur/flock-api/app/models"
	"net/http"
)

func PagesIndex(w http.ResponseWriter, r *http.Request) {
	page := new(PagesAttributes)
	page.Name = "home"

	users := models.Users{
		models.User{Name: "Omar", Email: "oozgur217@gmail.com"},
		models.User{Name: "Aditya", Email: "rajuaditya@gmail.com"},
	}

	json.NewEncoder(w).Encode(users)
}

type PagesAttributes struct {
	Name string
}
