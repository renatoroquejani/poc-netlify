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
        "/api/deploy/site": {
            "post": {
                "description": "Cria um novo site na Netlify ou atualiza um existente quando o ID é fornecido",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "deploy"
                ],
                "summary": "Cria ou atualiza sites na Netlify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do site na Netlify para atualização (opcional)",
                        "name": "site_id",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Nome do site para teste",
                        "name": "site_name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Descrição do site para teste",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "boolean",
                        "description": "Remover o site após o teste",
                        "name": "cleanup_after",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Domínio personalizado para o site (opcional)",
                        "name": "custom_domain",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Arquivo HTML para deploy (opcional)",
                        "name": "file",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Caminho da pasta local para deploy (opcional)",
                        "name": "folder_path",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.TestDeployResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.TestDeployResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.TestDeployResponse"
                        }
                    }
                }
            }
        },
        "/api/domains/add": {
            "post": {
                "description": "Adiciona um domínio personalizado como alias para um site existente na Netlify",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "domains"
                ],
                "summary": "Adiciona um domínio personalizado a um site",
                "parameters": [
                    {
                        "description": "Dados do domínio a ser adicionado",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.DomainRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    }
                }
            }
        },
        "/api/domains/remove": {
            "post": {
                "description": "Remove um domínio personalizado dos aliases de um site existente na Netlify",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "domains"
                ],
                "summary": "Remove um domínio personalizado de um site",
                "parameters": [
                    {
                        "description": "Dados do domínio a ser removido",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.DomainRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    }
                }
            }
        },
        "/api/domains/remove-primary": {
            "post": {
                "description": "Remove o domínio principal de um site existente na Netlify, mantendo os aliases",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "domains"
                ],
                "summary": "Remove o domínio principal de um site",
                "parameters": [
                    {
                        "description": "ID do site a ter o domínio principal removido",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.DomainRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    }
                }
            }
        },
        "/api/domains/set-default": {
            "post": {
                "description": "Define um domínio personalizado como o domínio principal de um site existente na Netlify",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "domains"
                ],
                "summary": "Define um domínio como o domínio principal",
                "parameters": [
                    {
                        "description": "Dados do domínio a ser definido como principal",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.DomainRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.DomainResponse"
                        }
                    }
                }
            }
        },
        "/api/test/netlify/connection": {
            "get": {
                "description": "Testa a conexão com a API da Netlify e exibe informações sobre o token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "netlify"
                ],
                "summary": "Testa a conexão com a API da Netlify",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.DomainRequest": {
            "type": "object",
            "required": [
                "site_id"
            ],
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "example.com"
                },
                "site_id": {
                    "type": "string",
                    "example": "e17e2166-d8ab-4cad-9916-a9a3fed7750d"
                }
            }
        },
        "api.DomainResponse": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "example.com"
                },
                "message": {
                    "type": "string",
                    "example": "Domínio adicionado com sucesso"
                },
                "site_id": {
                    "type": "string",
                    "example": "e17e2166-d8ab-4cad-9916-a9a3fed7750d"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "api.TestDeployResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2023-10-25T15:30:45Z"
                },
                "message": {
                    "type": "string",
                    "example": "Site de teste criado com sucesso"
                },
                "site_id": {
                    "type": "string",
                    "example": "a1b2c3d4"
                },
                "site_url": {
                    "type": "string",
                    "example": "https://test-site.netlify.app"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "test_success": {
                    "type": "boolean",
                    "example": true
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Netlify Deploy API",
	Description:      "API para facilitar o deploy de arquivos estáticos na Netlify a partir de um bucket S3",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
