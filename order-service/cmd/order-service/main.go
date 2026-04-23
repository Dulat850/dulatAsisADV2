package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    
    "dulatAsisADV2/order-service/internal/config"
    "dulatAsisADV2/order-service/internal/repository"
    "dulatAsisADV2/order-service/internal/transport/grpc"
    "dulatAsisADV2/order-service/internal/transport/http"
    "dulatAsisADV2/order-service/internal/usecase"
)

func main() {
    cfg := config.Load()
    
    // Database connection
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    log.Println("Database connected successfully")
    
    // Initialize repository
    orderRepo := repository.NewPostgresOrderRepo(db)
    
    // Initialize use case
    orderUseCase := usecase.NewOrderUseCase(orderRepo)
    
    // Start gRPC server in goroutine
    go func() {
        if err := grpc.StartGRPCServer(cfg.GRPCPort, orderUseCase); err != nil {
            log.Fatal("Failed to start gRPC server:", err)
        }
    }()
    
    // Start HTTP server
    router := http.NewRouter(orderUseCase)
    log.Printf("Order Service HTTP server listening on port %s", cfg.HTTPPort)
    if err := router.Run(":" + cfg.HTTPPort); err != nil {
        log.Fatal("Failed to start HTTP server:", err)
    }
}
