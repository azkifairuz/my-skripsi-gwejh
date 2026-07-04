package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/database"
	financev1 "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/generated/proto/finance/v1"
	grpcserver "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/grpc"
	"github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	server := grpc.NewServer()
	transactionRepository := setupTransactionRepository()
	financev1.RegisterFinanceServiceServer(server, grpcserver.NewServer(transactionRepository))
	reflection.Register(server)

	log.Printf("finance gRPC server listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func setupTransactionRepository() *repository.TransactionRepository {
	config := database.LoadConfig()
	if !config.Enabled() {
		log.Print("database is not configured; CreateTransaction will be unavailable")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Open(ctx, config)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Print("database connected and migrated")
	return repository.NewTransactionRepository(db)
}
