package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "project-manager/docs"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/api/option"
)

// @title Project Manager API
// @version 1.0
// @description API for managing projects
// @host localhost:8080
// @BasePath /api

// Project struct
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

// @title Packages Manager API
// @version 1.0
// @description API for managing packages
// @host localhost:8080
// @BasePath /api

// Packages struct
type Packages struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func InitFirebase() (*firestore.Client, error) {
	println("Initializing Firebase with Railway shared variables")

	// Retrieve the Firebase credentials JSON from the environment variable
	credentialsJSON := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsJSON == "" {
		return nil, fmt.Errorf("firebase credentials not found in environment variables")
	}

	// Use the credentials to initialize Firebase
	sa := option.WithCredentialsJSON([]byte(credentialsJSON))
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, err
	}

	// Connect to Firestore
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateProjectHandler godoc
// @Summary Create a new project
// @Description Create a new project with the provided data
// @Tags projects
// @Accept  json
// @Produce  json
// @Param project body Project true "Project data"
// @Success 201 {object} Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/new [post]
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project Project

	// Parse JSON data
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if project.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	// Add project to Firestore
	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	docRef, _, err := client.Collection("projects-manager").Add(context.Background(), project)
	if err != nil {
		http.Error(w, "Error saving project to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the generated Firestore ID into the project struct
	project.ID = docRef.ID

	// Return the project including the ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// Get all projects (GET)
// @Summary Get all projects
// @Description Retrieve a list of all projects
// @Produce json
// @Success 200 {array} Project
// @Router /projects [get]
func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {

	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	var projects []Project
	iter := client.Collection("projects-manager").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var project Project
		doc.DataTo(&project)
		project.ID = doc.Ref.ID // Set the project ID
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Get single project by ID (GET)
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} Project
// @Router /projects/{id} [get]
func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	doc, err := client.Collection("projects-manager").Doc(projectID).Get(context.Background())
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var project Project
	doc.DataTo(&project)
	project.ID = doc.Ref.ID // Set the project ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Update a project by ID (PUT)
// @Summary Update a project by ID
// @Description Update an existing project by its ID
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Project ID"
// @Param name formData string false "Project name"
// @Param imageUrl formData string false "Image URL"
// @Param link formData string false "Project link"
// @Param description formData string false "Project description"
// @Success 200 {object} map[string]string
// @Router /projects/{id} [put]
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Fetch the existing project from Firebase
	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	projectRef := client.Collection("projects-manager").Doc(projectID)
	doc, err := projectRef.Get(context.Background())
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Get the existing project data
	var existingProject Project
	doc.DataTo(&existingProject)

	// Update only the fields provided in the form
	if name := r.FormValue("name"); name != "" {
		existingProject.Name = name
	}
	if imageUrl := r.FormValue("imageUrl"); imageUrl != "" {
		existingProject.ImageUrl = imageUrl
	}
	if link := r.FormValue("link"); link != "" {
		existingProject.Link = link
	}
	if description := r.FormValue("description"); description != "" {
		existingProject.Description = description
	}

	// Save the updated project
	_, err = projectRef.Set(context.Background(), existingProject)
	if err != nil {
		http.Error(w, "Error updating project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project updated successfully"})
}

// Delete a project by ID (DELETE)
// @Summary Delete a project by ID
// @Description Delete a project by its ID
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]string
// @Router /projects/{id} [delete]
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Delete the project
	_, err = client.Collection("projects-manager").Doc(projectID).Delete(context.Background())
	if err != nil {
		http.Error(w, "Error deleting project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
}

// CreatePackageHandler godoc
// @Summary Create a new package
// @Description Create a new package with the provided data
// @Tags packages
// @Accept  json
// @Produce  json
// @Param package body Packages true "Packages data"
// @Success 201 {object} Packages
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/packages/new [post]
func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	var packages Packages

	// Parse JSON data
	if err := json.NewDecoder(r.Body).Decode(&packages); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if packages.Name == "" {
		http.Error(w, "Packages name is required", http.StatusBadRequest)
		return
	}

	// Add packages to Firestore
	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	docRef, _, err := client.Collection("packages-manager").Add(context.Background(), packages)
	if err != nil {
		http.Error(w, "Error saving packages to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the generated Firestore ID into the packages struct
	packages.ID = docRef.ID

	// Return the packages including the ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(packages)
}

// Get all packages (GET)
// @Summary Get all packages
// @Description Retrieve a list of all packages
// @Produce json
// @Success 200 {array} Packages
// @Router /packages [get]
func GetPackagesHandler(w http.ResponseWriter, r *http.Request) {
	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	var packages []Packages
	iter := client.Collection("packages-manager").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var packageObj Packages
		doc.DataTo(&packageObj)
		packageObj.ID = doc.Ref.ID // Set the packages ID
		packages = append(packages, packageObj)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

// Get single package by ID (GET)
// @Summary Get a package by ID
// @Description Retrieve a package by its ID
// @Produce json
// @Param id path string true "Packages ID"
// @Success 200 {object} Packages
// @Router /packages/{id} [get]
func GetPackageByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	doc, err := client.Collection("packages-manager").Doc(projectID).Get(context.Background())
	if err != nil {
		http.Error(w, "Packages not found", http.StatusNotFound)
		return
	}

	var packages Packages
	doc.DataTo(&packages)
	packages.ID = doc.Ref.ID // Set the package ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

// Update a package by ID (PUT)
// @Summary Update a package by ID
// @Description Update an existing package by its ID
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Packages ID"
// @Param name formData string false "Packages name"
// @Param imageUrl formData string false "Image URL"
// @Param link formData string false "Packages link"
// @Param description formData string false "Packages description"
// @Success 200 {object} map[string]string
// @Router /packages/{id} [put]
func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Fetch the existing package from Firebase
	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	projectRef := client.Collection("packages-manager").Doc(projectID)
	doc, err := projectRef.Get(context.Background())
	if err != nil {
		http.Error(w, "Packages not found", http.StatusNotFound)
		return
	}

	// Get the existing package data
	var existingProject Packages
	doc.DataTo(&existingProject)

	// Update only the fields provided in the form
	if name := r.FormValue("name"); name != "" {
		existingProject.Name = name
	}

	if link := r.FormValue("link"); link != "" {
		existingProject.Link = link
	}
	if description := r.FormValue("description"); description != "" {
		existingProject.Description = description
	}

	// Save the updated package
	_, err = projectRef.Set(context.Background(), existingProject)
	if err != nil {
		http.Error(w, "Error updating package", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Packages updated successfully"})
}

// Delete a package by ID (DELETE)
// @Summary Delete a package by ID
// @Description Delete a package by its ID
// @Param id path string true "Packages ID"
// @Success 200 {object} map[string]string
// @Router /packages/{id} [delete]
func DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	client, err := InitFirebase()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Delete the package
	_, err = client.Collection("packages-manager").Doc(projectID).Delete(context.Background())
	if err != nil {
		http.Error(w, "Error deleting package", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Packages deleted successfully"})
}

func main() {

	// log.Println("FIREBASE_CREDENTIALS:", os.Getenv("FIREBASE_CREDENTIALS"))

	// Check if running in production or not
	if os.Getenv("RAILWAY_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	r := mux.NewRouter()
	// Define CRUD routes projects
	r.HandleFunc("/api/projects", GetProjectsHandler).Methods("GET")           // Get all projects
	r.HandleFunc("/api/projects/{id}", GetProjectByIDHandler).Methods("GET")   // Get a project by ID
	r.HandleFunc("/api/projects/new", CreateProjectHandler).Methods("POST")    // Create a new project
	r.HandleFunc("/api/projects/{id}", UpdateProjectHandler).Methods("PUT")    // Update a project
	r.HandleFunc("/api/projects/{id}", DeleteProjectHandler).Methods("DELETE") // Delete a project
	// Define CRUD routes for packages
	r.HandleFunc("/api/packages", GetPackagesHandler).Methods("GET")           // Get all packages
	r.HandleFunc("/api/packages/{id}", GetPackageByIDHandler).Methods("GET")   // Get a package by ID
	r.HandleFunc("/api/packages/new", CreatePackageHandler).Methods("POST")    // Create a new package
	r.HandleFunc("/api/packages/{id}", UpdatePackageHandler).Methods("PUT")    // Update a package
	r.HandleFunc("/api/packages/{id}", DeletePackageHandler).Methods("DELETE") // Delete a package

	// Serve Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Set up CORS options
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://project-manager-production-7def.up.railway.app"}), // Allow your frontend domain
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),                                                    // Allow HTTP methods
		handlers.AllowedHeaders([]string{"Content-Type"}),                                                                    // Allow necessary headers
	)

	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsOptions(r)))
}
