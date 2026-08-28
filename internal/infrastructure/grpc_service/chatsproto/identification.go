package chatsproto

import (
	"chats-service/configs"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type TokenSubject struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func GetTokenFromString(tokenString string) (*jwt.Token, error) {
	configuration := configs.GetAPIConfiguration()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(configuration.SecretKey), nil
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
