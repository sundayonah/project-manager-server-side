package models

// ProjectData represents the structure for creating or updating a project
type ProjectData struct {
	Name        string   `json:"name"`
	ImageUrl    string   `json:"imageUrl"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Stacks      []string `json:"stacks"` // Array of technology stacks
}

// PackageData represents the structure for creating or updating a package
type PackageData struct {
	Name        string   `json:"name"`
	Link        string   `json:"link,omitempty"`
	Description string   `json:"description,omitempty"`
	Stacks      []string `json:"stacks"` // Array of technology stacks
}

// ClientData represents the structure for creating or updating a client
type ClientData struct {
	Name     string `json:"name"`
	Link     string `json:"link,omitempty"`
	ImageUrl string `json:"imageUrl"`
}

// ProjectResponse is used when returning project details including ID
type ProjectResponse struct {
	ID int `json:"id"`
	ProjectData
}

// PackageResponse is used when returning package details including ID
type PackageResponse struct {
	ID int `json:"id"`
	PackageData
}

// ClientResponse is used when returning client details including ID
type ClientResponse struct {
	ID int `json:"id"`
	ClientData
}
