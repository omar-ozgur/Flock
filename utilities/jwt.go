package utilities

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

// FLOCK_TOKEN_SECRET is the flock token secret
var FLOCK_TOKEN_SECRET = []byte(os.Getenv("FLOCK_TOKEN_SECRET"))

// GetClaims retrieves claims from a token
func GetClaims(tokenString string) map[string]interface{} {

	// Check if the token is empty
	if tokenString == "" {
		return nil
	}

	// Get token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return FLOCK_TOKEN_SECRET, nil
	})

	// Get claims
	claims := token.Claims.(jwt.MapClaims)

	return claims
}
