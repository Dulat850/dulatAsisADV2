package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    
    "dulatAsisADV2/payment-service/internal/config"
    "dulatAsisADV2/payment-service/internal/repository"
    "dulatAsisADV2/payment-service/internal/transport/grpc"
    "dulatAsisADV2/payment-service/internal/usecase"
)

func main() {
    cfg := config.Load()

    // Database connection
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    log.Printf("Connecting to database: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Проверка соединения
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    log.Println("Database connected successfully")

    // Initialize repository
    paymentRepo := repository.NewPostgresPaymentRepo(db)

    // Initialize use case
    paymentUseCase := usecase.NewPaymentUseCase(paymentRepo)

    // Start gRPC server
    if err := grpc.StartGRPCServer(cfg.GRPCServerPort, paymentUseCase); err != nil {
        log.Fatal("Failed to start gRPC server:", err)
    }
}
