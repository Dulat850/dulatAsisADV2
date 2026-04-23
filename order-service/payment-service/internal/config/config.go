package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBHost         string
    DBPort         int
    DBUser         string
    DBPassword     string
    DBName         string
    GRPCServerPort string
}

func Load() *Config {
    return &Config{
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         getEnvInt("DB_PORT", 5432),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", "postgres"),
        DBName:         getEnv("DB_NAME", "payment_db"),
        GRPCServerPort: getEnv("GRPC_PORT", "50052"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
