package controllers

import (
	"net/http"
	"golang.org/x/oauth2"
	"fmt"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"github.com/omar-ozgur/flock-api/utilities"
)

func LoginIndex(w http.ResponseWriter, r *http.Request) {
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

	page.Name = first_name + " " + last_name
	page.URL = tok.AccessToken
	templates.ExecuteTemplate(w, "test_login.html", page)
	
}

type LoggedInPageAttr struct {
	Name string
	URL string
}