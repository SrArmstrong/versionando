package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var claveSecreta = []byte("clave-super-secreta")

func CrearToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user": userID,
		"exp":  time.Now().Add(10 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(claveSecreta)
}

func ValidarToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return claveSecreta, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
