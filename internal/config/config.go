package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBDebug      string
	KafkaBrokers         string
	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
	MQTTBroker           string
}

// LoadConfig lee las variables de entorno desde el archivo .env si existe, o desde OS env vars
// Utiliza variables de entorno OS si existen, si no hace un fallback a `.env`, ideal para Prod.
func LoadConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No se pudo cargar el archivo env '%s', se usarán las variables del sistema (Modo Prod)", envFile)
	}

	return &Config{
		Port:         getEnv("PORT", "8000"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "ancianato"),
		DBDebug:      getEnv("DB_DEBUG", "false"),
		KafkaBrokers:         getEnv("KAFKA_BROKERS", "localhost:9092"),
		KeycloakURL:          getEnv("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM", "ancianato"),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID", "backend-client"),
		KeycloakClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		MQTTBroker:           getEnv("MQTT_BROKER", "tcp://localhost:1883"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
