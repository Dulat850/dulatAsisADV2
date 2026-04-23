package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
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
    // Load .env file if exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Parse DB port
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        dbPort = 5432
    }

    return &Config{
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         dbPort,
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", ""),
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
