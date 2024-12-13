package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000", "https://project-manager-server-side-production.up.railway.app"}, // Add your frontend URL here
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
			ExposedHeaders: []string{"Content-Length"},
			MaxAge:         86400,
		})

		return corsMiddleware.Handler(next)
	}
}
