package main

import (
	"log"
	"net"
	"os"

	ocrv1 "github.com/azkifairuz/my-skripsi-gwejh/ocr-service/internal/generated/proto/ocr/v1"
	grpcserver "github.com/azkifairuz/my-skripsi-gwejh/ocr-service/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	server := grpc.NewServer()
	ocrv1.RegisterOcrServiceServer(server, grpcserver.NewServer())
	reflection.Register(server)

	log.Printf("ocr gRPC server listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
