package main

import (
	"thanhldt060802/appconfig"
	"thanhldt060802/internal/opentelemetry"
	"thanhldt060802/internal/postgresqlclient"
	"thanhldt060802/internal/redisclient"
	"thanhldt060802/repository"
	"thanhldt060802/repository/db"
	"thanhldt060802/service"
)

func main() {
	appconfig.InitConfig()

	shutdown := opentelemetry.NewTracer(opentelemetry.TracerEndPointConfig{
		Host: appconfig.AppConfig.JaegerOTLPHost,
		Port: appconfig.AppConfig.JaegerOTLPPort,
	})
	defer shutdown()

	redisclient.RedisClient = redisclient.NewRedisClient(redisclient.RedisConfig{
		Host:     appconfig.AppConfig.RedisHost,
		Port:     appconfig.AppConfig.RedisPort,
		Database: appconfig.AppConfig.RedisDatabase,
		Password: appconfig.AppConfig.RedisPassword,
	})

	repository.BunSqlClient = postgresqlclient.NewBunSqlClient(postgresqlclient.BunSqlConfig{
		Host:     appconfig.AppConfig.PostgresHost,
		Port:     appconfig.AppConfig.PostgresPort,
		Database: appconfig.AppConfig.PostgresDatabase,
		Username: appconfig.AppConfig.PostgresUsername,
		Password: appconfig.AppConfig.PostgresPassword,
	})

	repository.PlayerRepo = db.NewPlayerRepo()

	playerService := service.NewPlayerService()
	playerService.InitSubscriber()

	select {}
}
