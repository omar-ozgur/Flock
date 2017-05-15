package controllers

import (
	"net/http"
	"golang.org/x/oauth2"
	"fmt"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"github.com/omar-ozgur/flock-api/utilities"
	"github.com/omar-ozgur/flock-api/app/models"
	"strconv"
)


func LoginWithFacebook(w http.ResponseWriter, r *http.Request) {
	page := new(LoggedInPageAttr)
	code := r.FormValue("code")
	tok, err := utilities.FbConfig.Exchange(oauth2.NoContext, code)
	fmt.Println(tok)

	if err != nil {
		fmt.Println(err)
	}

	response, err := http.Get("https://graph.facebook.com/me?access_token=" + tok.AccessToken + "&fields=email,first_name,last_name,id,friends")

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err == nil{
		//fmt.Println(body)
	}

	user, _ := jason.NewObjectFromBytes([]byte(body))
	fmt.Println(user)

	first_name, _ := user.GetString("first_name")
	last_name, _ := user.GetString("last_name")
	email, _:= user.GetString("email")
	fb_id_string, _:= user.GetString("id")
	fb_id, err := strconv.Atoi(fb_id_string)

	page.Name = first_name + " " + last_name

	app_token = models.ProcessFBLogin(first_name, last_name, email, fb_id)



	page.URL = tok.AccessToken
	


	templates.ExecuteTemplate(w, "test_login.html", page)
	
}

type LoggedInPageAttr struct {
	Name string
	URL string
}