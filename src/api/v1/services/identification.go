package services

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
)

func GetUserFromToken(tokenString string) (*protousers.UserResponse, error) {
	user, err := grpc_client.UsersGrpcClient.GetUserByToken(tokenString)

	if err != nil {
		log.Printf("Error when getting user by token: %v", err)
		return new(protousers.UserResponse), err
	}

	return user, nil
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *protousers.UserResponse
		authorization := r.Header["Authorization"]
		ctx := r.Context()

		if len(authorization) != 0 {
			token := strings.Replace(r.Header["Authorization"][0], "Bearer ", "", 1)
			user, _ = GetUserFromToken(token)
		}

		ctx = context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
