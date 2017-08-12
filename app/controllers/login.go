package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
)

func LoginWithFacebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, _ := ioutil.ReadAll(r.Body)
	JSONData, _ := jason.NewObjectFromBytes(b)
	token, _ := JSONData.GetString("token")

	response, _ := http.Get("https://graph.facebook.com/me?access_token=" + token + "&fields=email,first_name,last_name,id,friends")

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	user, _ := jason.NewObjectFromBytes([]byte(body))

	first_name, _ := user.GetString("first_name")
	last_name, _ := user.GetString("last_name")
	email, _ := user.GetString("email")
	fb_id_string, _ := user.GetString("id")

	status, message, app_token := models.ProcessFBLogin(first_name, last_name, email, fb_id_string)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":    status,
		"message":   message,
		"fb_token":  token,
		"app_token": app_token,
	})

	w.Write(JSON)

}

type LoggedInPageAttr struct {
	Name string
	URL  string
}
