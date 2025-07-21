package server

import (
	"fmt"
	"net/http"
	"thanhldt060802/appconfig"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewHTTPServer() *gin.Engine {
	engine := gin.New()
	engine.Use(otelgin.Middleware(appconfig.AppConfig.AppName))
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service-name": appconfig.AppConfig.AppName,
			"version":      appconfig.AppConfig.AppVersion,
			"time":         time.Now().Unix(),
		})
	})

	return engine
}

func Start(server *gin.Engine) {
	exit := make(chan struct{})
	go func() {
		if err := server.Run(fmt.Sprintf(":%v", appconfig.AppConfig.AppPort)); err != nil {
			log.Errorf("failed to start service %v: %v", appconfig.AppConfig.AppName, err.Error())
			close(exit)
		}
	}()
	log.Infof("service %v listening on port %v", appconfig.AppConfig.AppName, appconfig.AppConfig.AppPort)
	<-exit
}
