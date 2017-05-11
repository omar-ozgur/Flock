package utilities

import (
	"golang.org/x/oauth2"
)

var FbConfig = &oauth2.Config{
 		// ClientId: FBAppID(string), ClientSecret : FBSecret(string)
 		// Example - ClientId: "1234567890", ClientSecret: "red2drdff6e2321e51aedcc94e19c76ee"

 		ClientID:     "163745930820314", // change this to yours
 		ClientSecret: "ba7714f1c0ad1bc5b5d6b498f76ea4d1",
 		RedirectURL:  "http://localhost:3000/loginWithFB", // change this to your webserver adddress
 		Scopes:       []string{"email", "user_friends", "public_profile"},
 		Endpoint: oauth2.Endpoint{
 			AuthURL:  "https://www.facebook.com/dialog/oauth",
 			TokenURL: "https://graph.facebook.com/oauth/access_token",
 		},
 	}