package main

import (
	"fmt"
	"net/http"
	"thanhldt060802/internal/sqlclient"
	"thanhldt060802/middleware/auth"
	"thanhldt060802/repository"
	"thanhldt060802/repository/db"
	server "thanhldt060802/server/http"
	"thanhldt060802/service"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"

	apiV1 "thanhldt060802/api/v1"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	switch viper.GetString("db.driver") {
	case "postgres":
		{
			sqlclient.SqlClientConnInstance = sqlclient.NewSqlClient(sqlclient.SqlConfig{
				Host:     viper.GetString("db.host"),
				Port:     viper.GetInt("db.port"),
				Database: viper.GetString("db.database"),
				Username: viper.GetString("db.username"),
				Password: viper.GetString("db.password"),
			})
		}
	}

	server.APP_NAME = viper.get

	initRepository()
}

func main() {

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
					URL:         fmt.Sprintf("http://%v:%v", appconfig.AppConfig.AppHost, appconfig.AppConfig.AppPort),
					Description: "Local Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
			},
		},
		OpenAPIPath:   fmt.Sprintf("/%v/openapi", appconfig.AppConfig.AppName),
		DocsPath:      "",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	}

	router.GET(fmt.Sprintf("/%v/api-document", appconfig.AppConfig.AppName), func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!doctype html>
		<html>
			<head>
				<title>MyService APIs</title>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script id="api-reference" data-url="/`+appconfig.AppConfig.AppName+`/openapi.json"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
		</html>
		`))
	})

	humaAPI := humagin.New(router, humaConfig)
	api := hureg.NewAPIGen(humaAPI)
	api = api.AddBasePath(fmt.Sprintf("%v/%v", appconfig.AppConfig.AppName, appconfig.AppConfig.AppVersion[:2]))

	auth.AuthMdw = auth.NewSimpleAuthMiddleware()

	apiV1.RegisterAPITask(api, service.NewTaskService(repository.TaskRepo))

	server.Start(router)
}

func initRepository() {
	repository.TaskRepo = db.NewTaskRepo(sqlclient.SqlClientConnInstance.GetDB())
}
