package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type IConfig struct {
	ServerAddr string
	Mode       string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string

	RabbitUser       string
	RabbitPassword   string
	RabbitHost       string
	RabbitPort       string
	RabbitTaskQueue  string
	RabbitAgentQueue string
	RabbitDSN        string

	AgentTimeout     int
	AgentPing        int
	AgentResolveTime int

	JwtExpires int
	JwtSecret  string
}

var config *IConfig

func Init() (*IConfig, error) {
	var err error

	if err = godotenv.Load(); err != nil {
		return nil, err
	}

	config = &IConfig{
		ServerAddr: getEnv("SERVER_ADDR", "localhost:8080"),
		Mode:       getEnv("MODE", "debug"),

		PostgresUser:     getEnv("POSTGRES_USER", "calculator"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "test12345"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresDatabase: getEnv("POSTGRES_DB", "calculator"),

		RabbitUser:       getEnv("RABBIT_USER", "guest"),
		RabbitPassword:   getEnv("RABBIT_PASSWORD", "guest"),
		RabbitHost:       getEnv("RABBIT_HOST", "localhost"),
		RabbitPort:       getEnv("RABBIT_PORT", "5672"),
		RabbitTaskQueue:  getEnv("RABBIT_TASK_QUEUE", "CalculatorTaskQueue1"),
		RabbitAgentQueue: getEnv("RABBIT_AGENT_QUEUE", "CalculatorAgentQueue1"),

		AgentTimeout:     getEnvInt("AGENT_TIMEOUT", 600),
		AgentPing:        getEnvInt("AGENT_PING", 60),
		AgentResolveTime: getEnvInt("AGENT_RESOLVE_TIME", 600),

		JwtExpires: getEnvInt("JWT_EXPIRES", 86400),
		JwtSecret:  getEnv("JWT_SECRET", "Calcul@tor"),
	}

	// build a dsn connection string for RabbitMQ
	config.RabbitDSN = fmt.Sprintf("amqp://%s:%s@%s:%s", config.RabbitUser, config.RabbitPassword, config.RabbitHost, config.RabbitPort)

	return config, nil
}

func Get() *IConfig {
	return config
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	logrus.Errorf("value from .env for \"%s\" doesn't found, set default value: %s", key, defaultVal)
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			logrus.Errorf("value from .env for \"%s\" isn't integer, set default value: %d", key, defaultVal)
			return defaultVal
		}
		return valueInt
	}

	logrus.Errorf("value from .env for \"%s\" doesn't found, set default value: %d", key, defaultVal)
	return defaultVal
}
