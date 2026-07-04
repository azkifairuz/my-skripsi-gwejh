package main

import (
	"log"
	"net"
	"os"

	"github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/client"
	whatsappv1 "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/generated/proto/whatsapp/v1"
	grpcserver "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	server := grpc.NewServer()
	dependencies := client.NewDependencies(
		getEnv("FINANCE_SERVICE_ADDR", "finance-service:50053"),
		getEnv("NLP_SERVICE_ADDR", "nlp-service:50054"),
		getEnv("OCR_SERVICE_ADDR", "ocr-service:50051"),
	)
	whatsappv1.RegisterWhatsappServiceServer(server, grpcserver.NewServer(dependencies))
	reflection.Register(server)

	log.Printf("whatsapp gRPC server listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
