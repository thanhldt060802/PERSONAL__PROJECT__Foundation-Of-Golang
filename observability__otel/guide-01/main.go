package main

import (
	"net/http"
	"thanhldt060802/internal/otelclient"
	"thanhldt060802/middleware/auth"
	"thanhldt060802/repository"
	"thanhldt060802/repository/db"
	server "thanhldt060802/server/http"
	"thanhldt060802/service"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	apiV1 "thanhldt060802/api/v1"
)

func main() {
	shutdown := otelclient.InitTracer()
	defer shutdown()

	router := server.NewHTTPServer()

	humaConfig := huma.Config{
		OpenAPI: &huma.OpenAPI{
			Components: &huma.Components{
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"standard-auth": {
						Type:         "http",
						Scheme:       "bearer",
						In:           "header",
						Description:  "Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\"",
						BearerFormat: "Token String",
						Name:         "Authorization",
					},
				},
			},
			Servers: []*huma.Server{
				{
					URL:         "http://localhost:8000",
					Description: "Local Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
			},
		},
		OpenAPIPath:   "/my-guide/openapi",
		DocsPath:      "",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	}

	router.GET("/my-guide/api-document", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!doctype html>
		<html>
			<head>
				<title>MyGuide APIs</title>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script id="api-reference" data-url="/my-guide/openapi.json"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
		</html>
		`))
	})

	humaAPI := humagin.New(router, humaConfig)
	api := hureg.NewAPIGen(humaAPI)
	api = api.AddBasePath("my-guide/v1")

	auth.AuthMdw = auth.NewSimpleAuthMiddleware()

	repository.UserRepo = db.NewUserRepo()

	apiV1.RegisterAPIExample(api, service.NewUserService())

	server.Start(router, "8000")
}
