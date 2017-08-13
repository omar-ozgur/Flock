package utilities

import (
	"golang.org/x/oauth2"
	"os"
)

// Facebook configuration
var FbConfig = &oauth2.Config{
	ClientID:     os.Getenv("FLOCK_FB_CLIENT_ID"),
	ClientSecret: os.Getenv("FLOCK_FB_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("FLOCK_FB_REDIRECT_URL"),
	Scopes:       []string{"email", "user_friends", "public_profile"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.facebook.com/dialog/oauth",
		TokenURL: "https://graph.facebook.com/oauth/access_token",
	},
}
