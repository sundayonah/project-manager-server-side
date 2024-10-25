package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "project-manager/docs"
	"project-manager/ent"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

var client *ent.Client

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
// @Param project body Projects true "Projects data"
// @Success 201 {object} ent.Projects
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/new [post]
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var projectData struct {
		Name        string `json:"name"`
		ImageUrl    string `json:"imageUrl"`
		Link        string `json:"link"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&projectData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if projectData.Name == "" {
		http.Error(w, "Projects name is required", http.StatusBadRequest)
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

	// Create project using Ent
	project, err := client.Projects.Create().
		SetName(projectData.Name).
		SetImageURL(projectData.ImageUrl).
		SetLink(projectData.Link).
		SetDescription(projectData.Description).
		Save(context.Background())

	if err != nil {
		http.Error(w, "Error creating project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// GetProjectsHandler godoc
// @Summary Get all projects
// @Description Retrieve a list of all projects
// @Produce json
// @Success 200 {array} ent.Projects
// @Router /projects [get]
func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := client.Projects.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error fetching projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetProjectByIDHandler godoc
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Router /projects/{id} [get]// GetProjectByIDHandler godoc
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Produce json
// @Param id path int true "Projects ID"
// @Success 200 {object} ent.Projects
// @Router /projects/{id} [get]
func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	// Convert string to int64
	var id int64
	var err error
	fmt.Sscan(idStr, &id, &err)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Convert int64 to int
	projectId := int(id)

	project, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		http.Error(w, "Projects not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
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

	// Convert string to int64
	var id int64
	var err error
	if idStr == "" {
		http.Error(w, "Missing project ID", http.StatusBadRequest)
		return
	}
	fmt.Sscan(idStr, &id, &err)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Name        *string `json:"name"`
		ImageUrl    *string `json:"imageUrl"`
		Link        *string `json:"link"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	update := client.Projects.UpdateOneID(int(id))
	if updateData.Name != nil && *updateData.Name != "" {
		update.SetName(*updateData.Name)
	}
	if updateData.ImageUrl != nil && *updateData.ImageUrl != "" {
		update.SetImageURL(*updateData.ImageUrl)
	}
	if updateData.Link != nil && *updateData.Link != "" {
		update.SetLink(*updateData.Link)
	}
	if updateData.Description != nil && *updateData.Description != "" {
		update.SetDescription(*updateData.Description)
	}

	project, err := update.Save(context.Background())
	if err != nil {
		http.Error(w, "Error updating project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
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

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize database connection
	var err error
	client, err = InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer client.Close()

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/projects", GetProjectsHandler).Methods("GET")
	r.HandleFunc("/api/projects/{id}", GetProjectByIDHandler).Methods("GET")
	r.HandleFunc("/api/projects/new", CreateProjectHandler).Methods("POST")
	r.HandleFunc("/api/projects/{id}", UpdateProjectHandler).Methods("PUT")
	r.HandleFunc("/api/projects/{id}", DeleteProjectHandler).Methods("DELETE")

	// Swagger setup
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS setup
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://your-domain.com"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsOptions(r)))
}
