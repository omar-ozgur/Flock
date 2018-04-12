package utilities

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

// FLOCK_TOKEN_SECRET is the flock token secret
var FLOCK_TOKEN_SECRET = []byte(os.Getenv("FLOCK_TOKEN_SECRET"))

// CreateToken creates a token based on given claims
func CreateToken(claims map[string]interface{}) string {

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	// Add claims to the token
	for key, value := range claims {
		tokenClaims[key] = value
	}

	// Generate the token string
	tokenString, _ := token.SignedString(FLOCK_TOKEN_SECRET)

	return tokenString
}

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
