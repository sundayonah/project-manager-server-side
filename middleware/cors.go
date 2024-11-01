package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
			ExposedHeaders: []string{"Content-Length"},
			MaxAge:         86400,
		})

		return corsMiddleware.Handler(next)
	}
}
