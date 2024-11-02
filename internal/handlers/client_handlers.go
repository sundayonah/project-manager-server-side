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

func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var clientData models.ClientData

	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if clientData.Name == "" {
		http.Error(w, "CLient name is required", http.StatusBadRequest)
		return
	}
	if clientData.Link == "" {
		http.Error(w, "Link is required", http.StatusBadRequest)
		return
	}
	if clientData.ImageUrl == "" {
		http.Error(w, "Image URL is required", http.StatusBadRequest)
		return
	}

	// Create client
	client, err := database.Client.Clients.Create().
		SetName(clientData.Name).
		SetLink(clientData.Link).
		SetImageUrl(clientData.ImageUrl).
		Save(context.Background())

	if err != nil {
		http.Error(w, "Error creating client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ClientResponse{
		ID:         client.ID,
		ClientData: clientData,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	clients, err := database.Client.Clients.Query().All(context.Background())
	if err != nil {
		http.Error(w, "Error fetching clients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []models.ClientResponse
	for _, client := range clients {

		response = append(response, models.ClientResponse{
			ID: client.ID,
			ClientData: models.ClientData{
				Name:     client.Name,
				Link:     client.Link,
				ImageUrl: client.ImageUrl,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetClientByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, err := database.Client.Clients.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "CLient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := models.ClientResponse{
		ID: client.ID,
		ClientData: models.ClientData{
			Name:     client.Name,
			Link:     client.Link,
			ImageUrl: client.ImageUrl,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	var id int64
	if idStr == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}
	fmt.Sscan(idStr, &id)

	var clientData models.ClientData
	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	update := database.Client.Clients.UpdateOneID(int(id))
	if clientData.Name != "" {
		update.SetName(clientData.Name)
	}
	if clientData.Link != "" {
		update.SetLink(clientData.Link)
	}
	if clientData.ImageUrl != "" {
		update.SetImageUrl(clientData.ImageUrl)
	}

	client, err := update.Save(context.Background())
	if err != nil {
		http.Error(w, "Error updating client", http.StatusInternalServerError)
		return
	}

	response := models.ClientResponse{
		ID: client.ID,
		ClientData: models.ClientData{
			Name:     client.Name,
			Link:     client.Link,
			ImageUrl: client.ImageUrl,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	// Convert string to int64
	var id int64
	var err error
	if idStr == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscan(idStr, &id); err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = database.Client.Clients.DeleteOneID(int(id)).Exec(context.Background())
	if err != nil {
		http.Error(w, "Error deleting client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "CLient deleted successfully"})
}
