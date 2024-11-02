package database

import (
	"context"
	"fmt"
	"os"

	"project-manager/ent"

	_ "github.com/lib/pq"
)

var Client *ent.Client

// InitDB initializes the database connection
func InitDB() (*ent.Client, error) {
	fmt.Println("Starting database connection")
	// Load environment variables
	connectionString := os.Getenv("DATABASE_URL")
	fmt.Println("connection string: ", connectionString)
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

	// Set the global client
	Client = client

	return client, nil
}
