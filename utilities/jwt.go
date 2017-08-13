package utilities

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

// JWT configuration values
var FLOCK_TOKEN_SECRET = []byte(os.Getenv("FLOCK_TOKEN_SECRET"))

// Get claims from a JWT token
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
