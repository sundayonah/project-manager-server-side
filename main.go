package main

import (
	"log"
	"net/http"

	"project-manager/internal/database"
	handler "project-manager/internal/handlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize the database
	client, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer client.Close()

	// Create a new router
	r := mux.NewRouter()

	// Project routes
	r.HandleFunc("/api/projects/new", handler.CreateProjectHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/projects", handler.GetProjectsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", handler.GetProjectByIDHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", handler.UpdateProjectHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", handler.DeleteProjectHandler).Methods("DELETE", "OPTIONS")

	// Package routes
	r.HandleFunc("/api/packages/new", handler.CreatePackageHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/packages", handler.GetPackagesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", handler.GetPackageByIDHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", handler.UpdatePackageHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", handler.DeletePackageHandler).Methods("DELETE", "OPTIONS")

	// Swagger documentation route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Enhanced CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://project-manager-server-side-production.up.railway.app"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Content-Length"}),
		handlers.MaxAge(86400), // 24 hours
	)(r)

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
