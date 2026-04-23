package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DBHost            string
    DBPort            string
    DBUser            string
    DBPassword        string
    DBName            string
    HTTPPort          string
    GRPCPort          string
    PaymentServiceAddr string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    return &Config{
        DBHost:            getEnv("DB_HOST", "localhost"),
        DBPort:            getEnv("DB_PORT", "5432"),
        DBUser:            getEnv("DB_USER", "postgres"),
        DBPassword:        getEnv("DB_PASSWORD", ""),
        DBName:            getEnv("DB_NAME", "order_db"),
        HTTPPort:          getEnv("HTTP_PORT", "8080"),
        GRPCPort:          getEnv("GRPC_PORT", "50051"),
        PaymentServiceAddr: getEnv("PAYMENT_SERVICE_ADDR", "localhost:50052"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
