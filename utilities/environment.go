package utilities

import (
	"golang.org/x/oauth2"
	"os"
)

var FbConfig = &oauth2.Config{
	// ClientId: FBAppID(string), ClientSecret : FBSecret(string)
	// Example - ClientId: "1234567890", ClientSecret: "red2drdff6e2321e51aedcc94e19c76ee"

	ClientID:     os.Getenv("FLOCK_FB_CLIENT_ID"), // change this to yours
	ClientSecret: os.Getenv("FLOCK_FB_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("FLOCK_FB_REDIRECT_URL"), // change this to your webserver adddress
	Scopes:       []string{"email", "user_friends", "public_profile"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.facebook.com/dialog/oauth",
		TokenURL: "https://graph.facebook.com/oauth/access_token",
	},
}
