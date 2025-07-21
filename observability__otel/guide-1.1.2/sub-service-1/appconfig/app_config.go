package appconfig

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AppName    string
	AppVersion string
	AppHost    string
	AppPort    int

	JaegerOTLPHost string
	JaegerOTLPPort int
}

var AppConfig *Config

func InitConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Load file .env failed: ", err)
	}

	AppConfig = &Config{
		AppName:    GetString("APP_NAME", "my-service"),
		AppVersion: GetString("APP_VERSION", "v1.0.1"),
		AppHost:    GetString("APP_HOST", "localhost"),
		AppPort:    GetInt("APP_PORT", 8000),

		JaegerOTLPHost: GetString("JAEGER_OTLP_HOST", "localhost"),
		JaegerOTLPPort: GetInt("JAEGER_OTLP_PORT", 4318),
	}

	log.Info("Load .env file successful")
}

func GetString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}

func GetInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		convertedValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Fatal("Get int from file .env failed: ", err)
		}
		return int(convertedValue)
	} else {
		return defaultValue
	}
}
