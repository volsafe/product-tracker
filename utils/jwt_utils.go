package utils

import (
	"errors"
	"profile/config"

	"github.com/golang-jwt/jwt"
)

var c = config.GetConfig()

func ExtractUserIDFromToken(tokenString string) (uint, error) {

	secretKey := c.Jwt.Secret

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}
	return 0, errors.New("invalid token")
}
