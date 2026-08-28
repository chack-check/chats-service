package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func UserRequired(token *jwt.Token) error {
	if token == nil {
		return fmt.Errorf("incorrect token")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err == nil && token.Valid && exp.Unix() > time.Now().Unix() {
		return nil
	}

	return fmt.Errorf("incorrect token")
}
