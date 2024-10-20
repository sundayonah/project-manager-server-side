// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/packages/new": {
            "post": {
                "description": "Create a new package with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "packages"
                ],
                "summary": "Create a new package",
                "parameters": [
                    {
                        "description": "Packages data",
                        "name": "package",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Packages"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Packages"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/projects/new": {
            "post": {
                "description": "Create a new project with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Create a new project",
                "parameters": [
                    {
                        "description": "Project data",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Project"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/packages": {
            "get": {
                "description": "Retrieve a list of all packages",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all packages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Packages"
                            }
                        }
                    }
                }
            }
        },
        "/packages/{id}": {
            "get": {
                "description": "Retrieve a package by its ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get a package by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Packages ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Packages"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing package by its ID",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update a package by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Packages ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Packages name",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Image URL",
                        "name": "imageUrl",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Packages link",
                        "name": "link",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Packages description",
                        "name": "description",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a package by its ID",
                "summary": "Delete a package by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Packages ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects": {
            "get": {
                "description": "Retrieve a list of all projects",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all projects",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Project"
                            }
                        }
                    }
                }
            }
        },
        "/projects/{id}": {
            "get": {
                "description": "Retrieve a project by its ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get a project by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Project"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing project by its ID",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update a project by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Project name",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Image URL",
                        "name": "imageUrl",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Project link",
                        "name": "link",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Project description",
                        "name": "description",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a project by its ID",
                "summary": "Delete a project by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Packages": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.Project": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "imageUrl": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Packages Manager API",
	Description:      "API for managing packages",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
