package middlewares

import (
	"net/http"

	"github.com/chack-check/chats-service/settings"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", settings.Settings.APP_ALLOW_ORIGINS)
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}
