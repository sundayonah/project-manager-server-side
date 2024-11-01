package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"project-manager/ent"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

var client *ent.Client

// Update the project data structure
type ProjectData struct {
	Name        string   `json:"name"`
	ImageUrl    string   `json:"imageUrl"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Stacks      []string `json:"stacks"` // Array of technology stacks
}

// Update the package data structure
type PackageData struct {
	Name        string   `json:"name"`
	Link        string   `json:"link,omitempty"`
	Description string   `json:"description,omitempty"`
	Stacks      []string `json:"stacks"` // Array of technology stacks
}

type ProjectResponse struct {
	ID int `json:"id"`
	ProjectData
}

type PackageResponse struct {
	ID int `json:"id"`
	PackageData
}

// Initialize Database connection
func InitDB() (*ent.Client, error) {
	// Load environment variables
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return nil, fmt.Errorf("database URL not found in environment variables")
	}

	// Open connection to PostgreSQL
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %v", err)
	}

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}

// CreateProjectHandler godoc
// @Summary Create a new project
// @Description Create a new project with the provided data
// @Tags projects
// @Accept json
// @Produce json
// @Param project body ent.Projects true "Projects data" // Updated reference
// @Success 201 {object} ent.Projects
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/new [post]

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var projectData ProjectData

	if err := json.NewDecoder(r.Body).Decode(&projectData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if projectData.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}
	if projectData.ImageUrl == "" {
		http.Error(w, "Image URL is required", http.StatusBadRequest)
		return
	}
	if projectData.Link == "" {
		http.Error(w, "Link is required", http.StatusBadRequest)
		return
	}
	if projectData.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	// Convert stacks array to JSON string
	stacksJSON, err := json.Marshal(projectData.Stacks)
	if err != nil {
		http.Error(w, "Error processing stacks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create project
	project, err := client.Projects.Create().
		SetName(projectData.Name).
		SetImageUrl(projectData.ImageUrl).
		SetLink(projectData.Link).
		SetDescription(projectData.Description).
		SetStacks(string(stacksJSON)).
		Save(context.Background())

	if err != nil {
		http.Error(w, "Error creating project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProjectResponse{
		ID:          project.ID,
		ProjectData: projectData,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetProjectsHandler godoc
// @Summary Get all projects
// @Description Retrieve a list of all projects
// @Produce json
// @Success 200 {array} ent.Projects
// @Router /api/projects [get]
func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := client.Projects.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error fetching projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []ProjectResponse
	for _, project := range projects {
		var stacks []string
		if err := json.Unmarshal([]byte(project.Stacks), &stacks); err != nil {
			stacks = []string{}
		}

		response = append(response, ProjectResponse{
			ID: project.ID,
			ProjectData: ProjectData{
				Name:        project.Name,
				ImageUrl:    project.ImageUrl,
				Link:        project.Link,
				Description: project.Description,
				Stacks:      stacks,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProjectByIDHandler godoc
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Router /projects/{id} [get]
func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := client.Projects.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	var stacks []string
	if err := json.Unmarshal([]byte(project.Stacks), &stacks); err != nil {
		stacks = []string{}
	}

	response := ProjectResponse{
		ID: project.ID,
		ProjectData: ProjectData{
			Name:        project.Name,
			ImageUrl:    project.ImageUrl,
			Link:        project.Link,
			Description: project.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProjectHandler godoc
// @Summary Update a project
// @Description Update a project by ID
// @Accept json
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Router /projects/{id} [put]
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	var id int64
	if idStr == "" {
		http.Error(w, "Missing project ID", http.StatusBadRequest)
		return
	}
	fmt.Sscan(idStr, &id)

	var projectData ProjectData
	if err := json.NewDecoder(r.Body).Decode(&projectData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	update := client.Projects.UpdateOneID(int(id))
	if projectData.Name != "" {
		update.SetName(projectData.Name)
	}
	if projectData.ImageUrl != "" {
		update.SetImageUrl(projectData.ImageUrl)
	}
	if projectData.Link != "" {
		update.SetLink(projectData.Link)
	}
	if projectData.Description != "" {
		update.SetDescription(projectData.Description)
	}
	if len(projectData.Stacks) > 0 {
		stacksJSON, err := json.Marshal(projectData.Stacks)
		if err != nil {
			http.Error(w, "Error processing stacks: "+err.Error(), http.StatusInternalServerError)
			return
		}
		update.SetStacks(string(stacksJSON))
	}

	project, err := update.Save(context.Background())
	if err != nil {
		http.Error(w, "Error updating project", http.StatusInternalServerError)
		return
	}

	var stacks []string
	if err := json.Unmarshal([]byte(project.Stacks), &stacks); err != nil {
		stacks = []string{}
	}

	response := ProjectResponse{
		ID: project.ID,
		ProjectData: ProjectData{
			Name:        project.Name,
			ImageUrl:    project.ImageUrl,
			Link:        project.Link,
			Description: project.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProjectHandler godoc
// @Summary Delete a project
// @Description Delete a project by ID
// @Param id path int true "Projects ID"
// @Success 200 {object} map[string]string
// @Router /projects/{id} [delete]
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	// Convert string to int64
	var id int64
	var err error
	if idStr == "" {
		http.Error(w, "Missing project ID", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscan(idStr, &id); err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	err = client.Projects.DeleteOneID(int(id)).Exec(context.Background())
	if err != nil {
		http.Error(w, "Error deleting project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
}

// CreatePackagesHandler godoc
// @Summary Create a new package
// @Description Create a new package with the provided data
// @Tags packages
// @Accept json
// @Produce json
// @Param package body ent.Packages true "Packages data" // Updated reference
// @Success 201 {object} ent.Packages
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/packages/new [post]

func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	var packageData PackageData

	if err := json.NewDecoder(r.Body).Decode(&packageData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if packageData.Name == "" {
		http.Error(w, "Package name is required", http.StatusBadRequest)
		return
	}

	stacksJSON, err := json.Marshal(packageData.Stacks)
	if err != nil {
		http.Error(w, "Error processing stacks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	packageRecord, err := client.Packages.Create().
		SetName(packageData.Name).
		SetLink(packageData.Link).
		SetDescription(packageData.Description).
		SetStacks(string(stacksJSON)).
		Save(context.Background())

	if err != nil {
		http.Error(w, "Error saving package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := PackageResponse{
		ID:          packageRecord.ID,
		PackageData: packageData,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetProjectsHandler godoc
// @Summary Get all packages
// @Description Retrieve a list of all packages
// @Produce json
// @Success 200 {array} ent.Packages
// @Router /api/packages [get]
func GetPackagesHandler(w http.ResponseWriter, r *http.Request) {
	packages, err := client.Packages.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error retrieving packages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []PackageResponse
	for _, pkg := range packages {
		var stacks []string
		if err := json.Unmarshal([]byte(pkg.Stacks), &stacks); err != nil {
			stacks = []string{}
		}

		response = append(response, PackageResponse{
			ID: pkg.ID,
			PackageData: PackageData{
				Name:        pkg.Name,
				Link:        pkg.Link,
				Description: pkg.Description,
				Stacks:      stacks,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPackageByIDHandler godoc
// @Summary Get a package by ID
// @Description Retrieve a package by its ID
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Summary Get a package by ID
// @Description Retrieve a package by its ID
// @Produce json
// @Param id path int true "Packages ID"
// @Success 200 {object} ent.Packages
// @Router /packages/{id} [get]
func GetPackageByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	pkg, err := client.Packages.Get(context.Background(), packageID)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "Package not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving package: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var stacks []string
	if err := json.Unmarshal([]byte(pkg.Stacks), &stacks); err != nil {
		stacks = []string{}
	}

	response := PackageResponse{
		ID: pkg.ID,
		PackageData: PackageData{
			Name:        pkg.Name,
			Link:        pkg.Link,
			Description: pkg.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdatePackageHandler godoc
// @Summary Update a package
// @Description Update a package by ID
// @Accept json
// @Produce json
// @Param id path int true "Packages ID"
// @Success 200 {object} ent.Packages
// @Router /packages/{id} [put]
func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	var packageData PackageData
	if err := json.NewDecoder(r.Body).Decode(&packageData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	update := client.Packages.UpdateOneID(packageID)
	if packageData.Name != "" {
		update.SetName(packageData.Name)
	}
	if packageData.Link != "" {
		update.SetLink(packageData.Link)
	}
	if packageData.Description != "" {
		update.SetDescription(packageData.Description)
	}
	if len(packageData.Stacks) > 0 {
		stacksJSON, err := json.Marshal(packageData.Stacks)
		if err != nil {
			http.Error(w, "Error processing stacks: "+err.Error(), http.StatusInternalServerError)
			return
		}
		update.SetStacks(string(stacksJSON))
	}

	pkg, err := update.Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "Package not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating package: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var stacks []string
	if err := json.Unmarshal([]byte(pkg.Stacks), &stacks); err != nil {
		stacks = []string{}
	}

	response := PackageResponse{
		ID: pkg.ID,
		PackageData: PackageData{
			Name:        pkg.Name,
			Link:        pkg.Link,
			Description: pkg.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProjectHandler godoc
// @Summary Delete a package
// @Description Delete a package by ID
// @Param id path int true "packages ID"
// @Success 200 {object} map[string]string
// @Router /packages/{id} [delete]
func DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageIDStr := params["id"]

	// Convert packageID from string to int
	var packageID int
	_, err := fmt.Sscan(packageIDStr, &packageID)
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	err = client.Packages.DeleteOneID(packageID).Exec(context.Background())
	if err != nil {
		http.Error(w, "Error deleting package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Package deleted successfully"})
}

func main() {
	// Initialize the database
	var err error
	client, err = InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer client.Close()

	// Create a new router
	r := mux.NewRouter()

	// Define your routes
	r.HandleFunc("/api/projects/new", CreateProjectHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/projects", GetProjectsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", GetProjectByIDHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", UpdateProjectHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/projects/{id}", DeleteProjectHandler).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/packages/new", CreatePackageHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/packages", GetPackagesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", GetPackageByIDHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", UpdatePackageHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/packages/{id}", DeletePackageHandler).Methods("DELETE", "OPTIONS")

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
