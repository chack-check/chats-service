package middlewares

import (
	"chats-service/configs"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	configuration := configs.GetAPIConfiguration()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", configuration.AllowOrigins)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}
