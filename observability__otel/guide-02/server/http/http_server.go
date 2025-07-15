package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	SERVICE_NAME = "my-guide-service"
	VERSION      = "v0.0.1"
)

func NewHTTPServer() *gin.Engine {
	engine := gin.New()
	engine.Use(otelgin.Middleware(SERVICE_NAME))
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service-name": SERVICE_NAME,
			"version":      VERSION,
			"time":         time.Now().Unix(),
		})
	})

	return engine
}

func Start(server *gin.Engine, port string) {
	exit := make(chan struct{})
	go func() {
		if err := server.Run(":" + port); err != nil {
			log.Errorf("failed to start service %v: %v", SERVICE_NAME, err.Error())
			close(exit)
		}
	}()
	log.Infof("service %v listening on port %v", SERVICE_NAME, port)
	<-exit
}
