package main

import (
	"fmt"
	"net/http"
	"thanhldt060802/internal/opentelemetry"
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

var HOST string = "localhost"
var PORT string = "8000"

func main() {
	shutdown := opentelemetry.NewTracer()
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
					URL:         fmt.Sprintf("http://%v:%v", HOST, PORT),
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

	repository.PlayerRepo = db.NewPlayerRepo()

	apiV1.RegisterAPIExample(api, service.NewPlayerService())

	server.Start(router, PORT)
}
