package main

import (
	"thanhldt060802/appconfig"
	"thanhldt060802/common/pubsub"
	"thanhldt060802/common/tracer"
	"thanhldt060802/internal/opentelemetry"
	"thanhldt060802/internal/redisclient"
	"thanhldt060802/internal/sqlclient"
	"thanhldt060802/repository"
	"thanhldt060802/repository/db"
	"thanhldt060802/service"
)

func main() {
	appconfig.InitConfig()

	opentelemetry.ShutdownTracer = opentelemetry.NewTracer(opentelemetry.TracerEndPointConfig{
		ServiceName: appconfig.AppConfig.AppName,
		Host:        appconfig.AppConfig.JaegerOTLPHost,
		Port:        appconfig.AppConfig.JaegerOTLPPort,
	})
	defer opentelemetry.ShutdownTracer()

	sqlclient.SqlClientConnInstance = sqlclient.NewSqlClient(sqlclient.SqlConfig{
		Host:     appconfig.AppConfig.PostgresHost,
		Port:     appconfig.AppConfig.PostgresPort,
		Database: appconfig.AppConfig.PostgresDatabase,
		Username: appconfig.AppConfig.PostgresUsername,
		Password: appconfig.AppConfig.PostgresPassword,
	})

	redisclient.RedisClientConnInstance = redisclient.NewRedisClient(redisclient.RedisConfig{
		Host:     appconfig.AppConfig.RedisHost,
		Port:     appconfig.AppConfig.RedisPort,
		Database: appconfig.AppConfig.RedisDatabase,
		Password: appconfig.AppConfig.RedisPassword,
	})
	pubsub.RedisSubInstance = pubsub.NewRedisSub[*tracer.MessageTracing](redisclient.RedisClientConnInstance.GetClient())

	repository.PlayerRepo = db.NewPlayerRepo()

	playerService := service.NewPlayerService()
	playerService.InitSubscriber()

	select {}
}
