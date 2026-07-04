package main

import (
	"log"
	"net"
	"os"

	financev1 "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/generated/proto/finance/v1"
	grpcserver "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/grpc"
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
	financev1.RegisterFinanceServiceServer(server, grpcserver.NewServer())
	reflection.Register(server)

	log.Printf("finance gRPC server listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
