package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/chack-check/chats-service/infrastructure/api/settings"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt/v5"
)

type TokenSubject struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func GetTokenFromString(tokenString string) (*jwt.Token, error) {
	log.Printf("Fetching token from string")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.Settings.APP_SECRET_KEY), nil
	})

	log.Printf("Fetched token = %+v. err = %v", token, err)
	return token, err
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header["Authorization"]
		log.Printf("Authorization header: %v", authorization)
		ctx := r.Context()

		if len(authorization) != 0 {
			tokenString := strings.Replace(r.Header["Authorization"][0], "Bearer ", "", 1)
			log.Printf("Parsing token: %s", tokenString)
			token, err := GetTokenFromString(tokenString)
			if err == nil && token.Valid {
				log.Printf("Successfully parsd token: %v", token)
				ctx = context.WithValue(r.Context(), "token", token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			sentry.CaptureException(err)
			log.Printf("Error validating token: %v", err)
		}

		log.Printf("Auahorization token is nil")
		ctx = context.WithValue(r.Context(), "token", nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTokenSubject(token *jwt.Token) (TokenSubject, error) {
	tokenSubject := TokenSubject{}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error parsing token subject: %v", token)
		sentry.CaptureException(err)
		return tokenSubject, err
	}

	err = json.Unmarshal([]byte(subject), &tokenSubject)
	if err != nil {
		log.Printf("Error parsing token subject: %v", token)
		sentry.CaptureException(err)
		return tokenSubject, err
	}

	log.Printf("Parsed token subject: %v", tokenSubject)

	return tokenSubject, nil
}
