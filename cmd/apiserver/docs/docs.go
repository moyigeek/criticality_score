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
        "/histories": {
            "get": {
                "description": "Get score histories by git link",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get score histories",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Git link",
                        "name": "link",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Skip count",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Take count",
                        "name": "take",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.PageDTO-model_ResultDTO"
                        }
                    }
                }
            }
        },
        "/query-with-pagination": {
            "get": {
                "description": "Query the database with pagination support",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Query with pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table Name",
                        "name": "tableName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page Size",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "Confidence",
                        "name": "confidence",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/rankings": {
            "get": {
                "description": "Get ranking results, optionally including all details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get ranking results",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Skip count",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Take count",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Include details",
                        "name": "detail",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.PageDTO-model_RankingResultDTO"
                        }
                    }
                }
            }
        },
        "/results": {
            "get": {
                "description": "Search score results by git link\nNOTE: All details are ignored, should use /results/:scoreid to get details\nNOTE: Maxium take count is 1000",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Search score results by git link",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search query",
                        "name": "q",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Skip count",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Take count",
                        "name": "take",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.PageDTO-model_ResultDTO"
                        }
                    }
                }
            }
        },
        "/results/{scoreid}": {
            "get": {
                "description": "Get score results, including all details by scoreid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get score results",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Score ID",
                        "name": "scoreid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ResultDTO"
                        }
                    }
                }
            }
        },
        "/search-packages": {
            "get": {
                "description": "Search for packages in the specified table that match the search query",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Search packages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table Name",
                        "name": "tableName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Search Query",
                        "name": "searchQuery",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page Size",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/update-gitlink": {
            "post": {
                "description": "Update the git link for a specified package",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update git link",
                "parameters": [
                    {
                        "description": "Update Git Link Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateGitLinkRequest"
                        }
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
        "controller.UpdateGitLinkRequest": {
            "type": "object",
            "required": [
                "linkConfidence",
                "newGitLink",
                "packageName",
                "tableName"
            ],
            "properties": {
                "linkConfidence": {
                    "type": "string"
                },
                "newGitLink": {
                    "type": "string"
                },
                "packageName": {
                    "type": "string"
                },
                "tableName": {
                    "type": "string"
                }
            }
        },
        "model.PageDTO-model_RankingResultDTO": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.RankingResultDTO"
                    }
                },
                "start": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.PageDTO-model_ResultDTO": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultDTO"
                    }
                },
                "start": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.RankingResultDTO": {
            "type": "object",
            "properties": {
                "distDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultDistDetailDTO"
                    }
                },
                "distroScore": {
                    "type": "number"
                },
                "gitDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultGitMetadataDTO"
                    }
                },
                "gitScore": {
                    "type": "number"
                },
                "langDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultLangDetailDTO"
                    }
                },
                "langScore": {
                    "type": "number"
                },
                "link": {
                    "type": "string"
                },
                "ranking": {
                    "type": "integer"
                },
                "score": {
                    "type": "number"
                },
                "scoreID": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                }
            }
        },
        "model.ResultDTO": {
            "type": "object",
            "properties": {
                "distDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultDistDetailDTO"
                    }
                },
                "distroScore": {
                    "type": "number"
                },
                "gitDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultGitMetadataDTO"
                    }
                },
                "gitScore": {
                    "type": "number"
                },
                "langDetail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ResultLangDetailDTO"
                    }
                },
                "langScore": {
                    "type": "number"
                },
                "link": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                },
                "scoreID": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                }
            }
        },
        "model.ResultDistDetailDTO": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "impact": {
                    "type": "number"
                },
                "pageRank": {
                    "type": "number"
                },
                "type": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                }
            }
        },
        "model.ResultGitMetadataDTO": {
            "type": "object",
            "properties": {
                "commitFrequency": {
                    "type": "number"
                },
                "contributorCount": {
                    "type": "integer"
                },
                "createdSince": {
                    "type": "string"
                },
                "language": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "license": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "orgCount": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                },
                "updatedSince": {
                    "type": "string"
                }
            }
        },
        "model.ResultLangDetailDTO": {
            "type": "object",
            "properties": {
                "depCount": {
                    "type": "integer"
                },
                "langEcoImpact": {
                    "type": "number"
                },
                "type": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
