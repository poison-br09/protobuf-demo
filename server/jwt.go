package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

// Generate a JWT token
func generateToken(agentId string) (string, error) {
	claims := jwt.MapClaims{
		"agent_id": agentId,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// validate a JWT Token
func validateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["agent_id"].(string), nil
}
