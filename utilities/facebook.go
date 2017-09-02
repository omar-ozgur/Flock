package utilities

import (
	"golang.org/x/oauth2"
	"os"
)

// FbConfig provides Facebook configuration information
var FbConfig = &oauth2.Config{

	// The Facebook client ID
	ClientID: os.Getenv("FLOCK_FB_CLIENT_ID"),

	// The Facebook client secret
	ClientSecret: os.Getenv("FLOCK_FB_CLIENT_SECRET"),

	// The Facebook redirect URL after login
	RedirectURL: os.Getenv("FLOCK_FB_REDIRECT_URL"),

	// The desired Facebook token scopes
	Scopes: []string{"email", "user_friends", "public_profile"},

	// The Facebook API endpoint
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.facebook.com/dialog/oauth",
		TokenURL: "https://graph.facebook.com/oauth/access_token",
	},
}
