package utilities

import (
	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secret")

func GetClaims(tokenString string) map[string]interface{} {
	if tokenString == "" {
		return nil
	}

	tokenString = tokenString[len("Bearer "):]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	return claims
}
