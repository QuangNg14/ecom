package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost            string
	Port                  string
	DBUser                string
	DBPassword            string
	DBAddress             string
	DBName                string
	JWTExpirationInSecond int64
	JWTSecret             string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:            getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                  getEnv("PORT", "8080"),
		DBUser:                getEnv("DB_USER", "root"),
		DBPassword:            getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:             fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                getEnv("DB_NAME", "ecom"),
		JWTExpirationInSecond: getEnvAsInt("JWT_EXPIRATION_IN_SECOND", 86400),
		JWTSecret:             getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return defaultVal
		}
		return i
	}
	return defaultVal
}
