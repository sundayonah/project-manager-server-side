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

// CreatePackageHandler handles the creation of a new package
func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	var packageData models.PackageData

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

	packageRecord, err := database.Client.Packages.Create().
		SetName(packageData.Name).
		SetLink(packageData.Link).
		SetDescription(packageData.Description).
		SetStacks(string(stacksJSON)).
		Save(context.Background())

	if err != nil {
		http.Error(w, "Error saving package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.PackageResponse{
		ID:          packageRecord.ID,
		PackageData: packageData,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetPackagesHandler retrieves all packages
func GetPackagesHandler(w http.ResponseWriter, r *http.Request) {
	packages, err := database.Client.Packages.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error retrieving packages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []models.PackageResponse
	for _, pkg := range packages {
		var stacks []string
		if err := json.Unmarshal([]byte(pkg.Stacks), &stacks); err != nil {
			stacks = []string{}
		}

		response = append(response, models.PackageResponse{
			ID: pkg.ID,
			PackageData: models.PackageData{
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

// GetPackageByIDHandler retrieves a package by its ID
func GetPackageByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	pkg, err := database.Client.Packages.Get(context.Background(), packageID)
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

	response := models.PackageResponse{
		ID: pkg.ID,
		PackageData: models.PackageData{
			Name:        pkg.Name,
			Link:        pkg.Link,
			Description: pkg.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	var packageData models.PackageData
	if err := json.NewDecoder(r.Body).Decode(&packageData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	update := database.Client.Packages.UpdateOneID(packageID)
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

	response := models.PackageResponse{
		ID: pkg.ID,
		PackageData: models.PackageData{
			Name:        pkg.Name,
			Link:        pkg.Link,
			Description: pkg.Description,
			Stacks:      stacks,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

	err = database.Client.Packages.DeleteOneID(packageID).Exec(context.Background())
	if err != nil {
		http.Error(w, "Error deleting package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Package deleted successfully"})
}
