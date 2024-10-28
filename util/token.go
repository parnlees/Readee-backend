package util

import (
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("yourSecretKey")

func GenerateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
