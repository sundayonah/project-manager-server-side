package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"project-manager/ent"
	"project-manager/internal/database"
	"project-manager/internal/models"

	"github.com/gorilla/mux"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var projectData models.ProjectData

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
	project, err := database.Client.Projects.Create().
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

	response := models.ProjectResponse{
		ID:          project.ID,
		ProjectData: projectData,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := database.Client.Projects.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error fetching projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []models.ProjectResponse
	for _, project := range projects {
		var stacks []string
		if err := json.Unmarshal([]byte(project.Stacks), &stacks); err != nil {
			stacks = []string{}
		}

		response = append(response, models.ProjectResponse{
			ID: project.ID,
			ProjectData: models.ProjectData{
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

func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := database.Client.Projects.Get(context.Background(), id)
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

	response := models.ProjectResponse{
		ID: project.ID,
		ProjectData: models.ProjectData{
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

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	var id int64
	if idStr == "" {
		http.Error(w, "Missing project ID", http.StatusBadRequest)
		return
	}
	fmt.Sscan(idStr, &id)

	var projectData models.ProjectData
	if err := json.NewDecoder(r.Body).Decode(&projectData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	update := database.Client.Projects.UpdateOneID(int(id))
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
			http.Error(w, "Error processing stacks", http.StatusInternalServerError)
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

	response := models.ProjectResponse{
		ID: project.ID,
		ProjectData: models.ProjectData{
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

	err = database.Client.Projects.DeleteOneID(int(id)).Exec(context.Background())
	if err != nil {
		http.Error(w, "Error deleting project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
}
