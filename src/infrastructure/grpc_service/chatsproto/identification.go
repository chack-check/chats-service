package chatsproto

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type TokenSubject struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func GetTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(Settings.APP_SECRET_KEY), nil
	})

	return token, err
}

func GetTokenSubject(token *jwt.Token) (TokenSubject, error) {
	tokenSubject := TokenSubject{}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return tokenSubject, err
	}

	err = json.Unmarshal([]byte(subject), &tokenSubject)
	if err != nil {
		return tokenSubject, err
	}

	return tokenSubject, nil
}
