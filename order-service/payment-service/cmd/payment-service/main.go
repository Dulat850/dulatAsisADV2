package main

import (
    "database/sql"
    "fmt"
    "log"
    "net"

    "dulatAsisADV2/payment-service/internal/config"
    "dulatAsisADV2/payment-service/internal/repository"
    grpcserver "dulatAsisADV2/payment-service/internal/transport/grpcserver"
    "dulatAsisADV2/payment-service/internal/usecase"

    _ "github.com/lib/pq"
    pb "dulatAsisADV2/proto/gen/payment"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.Load()

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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

    paymentRepo := repository.NewPostgresPaymentRepo(db)
    paymentUseCase := usecase.NewPaymentUseCase(paymentRepo)

    lis, err := net.Listen("tcp", ":"+cfg.GRPCServerPort)
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }

    grpcServer := grpc.NewServer()
    paymentServer := grpcserver.NewPaymentGRPCServer(paymentUseCase)
    pb.RegisterPaymentServiceServer(grpcServer, paymentServer)

    log.Printf("Payment Service gRPC server listening on port %s", cfg.GRPCServerPort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal("Failed to serve gRPC:", err)
    }
}
