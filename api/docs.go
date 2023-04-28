// Package api GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Ariel Gerardo Ríos",
            "url": "http://ariel17.com.ar/",
            "email": "arielgerardorios@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/ariel17/jobberwocky/blob/master/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/jobs": {
            "get": {
                "description": "Based on filter parameters it searches in jobs in the local database and in external resources concurrently.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "Search for published jobs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filters jobs by matching text in title or description (case-insensitive).",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filters jobs by matching company (case-sensitive).",
                        "name": "company",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filters jobs by matching location (case-sensitive).",
                        "name": "location",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filters jobs by matching salary, fixed or in range.",
                        "name": "salary",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filters jobs by matching work type (case-sensitive). Values: Full-Time, Contractor, Part-Time.",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filters jobs by remote condition.",
                        "name": "is_remote_friendly",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "Filters jobs by keywords (case-sensitive, inclusive).",
                        "name": "keywords",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Job"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new job receiving a JSON body with the details. If matching subscriptions exists, it sends notifications by email asynchronously.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "Publish a new job",
                "parameters": [
                    {
                        "description": "New job details.",
                        "name": "job",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Job"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Job"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/subscriptions": {
            "post": {
                "description": "Receives a JSON body with the email and filter values to match new job posts and be notified.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "Creates a new subscripion",
                "parameters": [
                    {
                        "description": "New subscription details.",
                        "name": "subscription",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Subscription"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Subscription"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Job": {
            "type": "object",
            "properties": {
                "company": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "is_remote_friendly": {
                    "type": "boolean"
                },
                "keywords": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "location": {
                    "type": "string"
                },
                "salary_max": {
                    "type": "integer"
                },
                "salary_min": {
                    "type": "integer"
                },
                "source": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "domain.Subscription": {
            "type": "object",
            "properties": {
                "company": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "is_remote_friendly": {
                    "type": "boolean"
                },
                "keywords": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "location": {
                    "type": "string"
                },
                "salary": {
                    "type": "integer"
                },
                "text": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "http.ErrorResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Jobberwocky API",
	Description:      "A job posting and searching API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
