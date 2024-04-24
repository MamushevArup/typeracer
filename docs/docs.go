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
        "/admin/auth/refresh": {
            "post": {
                "description": "Admin can refresh their token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "refresh for admins",
                "operationId": "refresh_admin",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.AdminSignInResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/admin/auth/sign-in": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "This endpoint is used for user authentication.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Sign in for admins",
                "operationId": "sign-in-admin",
                "parameters": [
                    {
                        "description": "Sign In",
                        "name": "models.AdminSignIn",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AdminSignIn"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SignInHandler"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/admin/moderation/all": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "This endpoint is used for see content in moderation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "moderation"
                ],
                "summary": "Get all text to see in pending mode",
                "operationId": "moderation-all",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ModerationServiceResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/admin/moderation/content/{moderation_id}/approve": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Admin can approve this content and it appear in global text storage",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "moderation"
                ],
                "summary": "Approve provided text content",
                "operationId": "moderation-approve",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Moderation ID",
                        "name": "moderation_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/admin/moderation/content/{moderation_id}/reject": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Admin can reject some content because of problem in content",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "moderation"
                ],
                "summary": "Reject provided content",
                "operationId": "moderation-reject",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Moderation ID",
                        "name": "moderation_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Reject",
                        "name": "models.ModerationRejectToService",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ModerationRejectToService"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/admin/moderation/{moderation_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get details of a specific moderation item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "moderation"
                ],
                "summary": "Moderation details",
                "operationId": "moderation-content",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Moderation ID",
                        "name": "moderation_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ModerationTextDetails"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/logout": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint is used for user logout.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Log out",
                "operationId": "log-out",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/refresh": {
            "post": {
                "description": "This endpoint is used to refresh the endpoint token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh token",
                "operationId": "refresh",
                "parameters": [
                    {
                        "description": "Refresh",
                        "name": "models.RefreshS",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RefreshS"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/sign-in": {
            "post": {
                "description": "This endpoint is used for user authentication.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in",
                "operationId": "sign-in",
                "parameters": [
                    {
                        "description": "Sign In",
                        "name": "models.SignIn",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignIn"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SignInHandler"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/sign-up": {
            "post": {
                "description": "This endpoint is used for user registration.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up",
                "operationId": "sign-up",
                "parameters": [
                    {
                        "description": "Sign Up",
                        "name": "models.SignUp",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignUp"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SignUpHandler"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/content/contribute": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Endpoint related to contribute text to the general text set",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "content"
                ],
                "summary": "Contribute text",
                "operationId": "contribute",
                "parameters": [
                    {
                        "description": "Contribute",
                        "name": "models.ContributeHandlerRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ContributeHandlerRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/single/curr-wpm": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint is used to calculate the current WPM for a racer.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "single"
                ],
                "summary": "Calculate current Words Per Minute (WPM)",
                "operationId": "curr-wpm",
                "parameters": [
                    {
                        "description": "Wpm calculation",
                        "name": "models.CountWpm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CountWpm"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Speed"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/single/end-race": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint is used to end a race for a racer.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "single"
                ],
                "summary": "End a race",
                "operationId": "end-race",
                "parameters": [
                    {
                        "description": "End Race",
                        "name": "models.ReqEndSingle",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ReqEndSingle"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.RespEndSingle"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/single/race": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint is used to start a new race for a racer.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "single"
                ],
                "summary": "Start a new race",
                "operationId": "start-race",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SingleResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/track/link": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "This endpoint is used to create a racetrack. It generates a unique link for the racetrack and returns it to the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "multiple"
                ],
                "summary": "Create a racetrack",
                "operationId": "create-racetrack",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.LinkCreation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/track/race/{link}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "This endpoint is used to join a racetrack. It upgrades the HTTP connection to a WebSocket connection. The server sends messages with the current race status to the client over the WebSocket connection.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "multiple"
                ],
                "summary": "Join a racetrack",
                "operationId": "racetrack",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Race Link",
                        "name": "link",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.RacerM"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.errorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.AdminSignIn": {
            "type": "object",
            "properties": {
                "fingerprint": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.AdminSignInResponse": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                }
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                }
            }
        },
        "models.ContributeHandlerRequest": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "source_title": {
                    "type": "string"
                }
            }
        },
        "models.CountWpm": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "index": {
                    "type": "integer"
                }
            }
        },
        "models.LinkCreation": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                }
            }
        },
        "models.ModerationRejectToService": {
            "type": "object",
            "properties": {
                "moderationID": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "models.ModerationServiceResponse": {
            "type": "object",
            "properties": {
                "contributor_name": {
                    "type": "string"
                },
                "moderation_id": {
                    "type": "string"
                },
                "sent_at": {
                    "type": "string"
                }
            }
        },
        "models.ModerationTextDetails": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "moderation_id": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "source_title": {
                    "type": "string"
                }
            }
        },
        "models.RacerInfo": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.RacerM": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.RefreshS": {
            "type": "object",
            "properties": {
                "fingerprint": {
                    "type": "string"
                }
            }
        },
        "models.ReqEndSingle": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "errors": {
                    "type": "integer"
                },
                "length": {
                    "type": "integer"
                }
            }
        },
        "models.RespEndSingle": {
            "type": "object",
            "properties": {
                "accuracy": {
                    "type": "number"
                },
                "duration": {
                    "type": "integer"
                },
                "wpm": {
                    "type": "integer"
                }
            }
        },
        "models.SignIn": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fingerprint": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.SignInHandler": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                },
                "avatar": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.SignUp": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fingerprint": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.SignUpHandler": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                },
                "avatar": {
                    "type": "string"
                }
            }
        },
        "models.SingleResponse": {
            "type": "object",
            "properties": {
                "racer": {
                    "$ref": "#/definitions/models.RacerInfo"
                },
                "text": {
                    "$ref": "#/definitions/models.TextInfo"
                }
            }
        },
        "models.Speed": {
            "type": "object",
            "properties": {
                "wpm": {
                    "type": "integer"
                }
            }
        },
        "models.TextInfo": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "content is actual text",
                    "type": "string"
                },
                "contributor_name": {
                    "description": "who contributed the content",
                    "type": "string"
                },
                "header": {
                    "description": "header the title of the text. Ex name of the book, song",
                    "type": "string"
                },
                "source": {
                    "description": "source from what it is coming. It can be from a book, article, etc.",
                    "type": "string"
                },
                "text_author": {
                    "description": "who wrote the content",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0",
	Host:             "localhost:1001",
	BasePath:         "",
	Schemes:          []string{"http"},
	Title:            "Typeracer Game clone API",
	Description:      "API for typeracer game clone. Typeracer popular game where users improve their typing skills in interactive format",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
