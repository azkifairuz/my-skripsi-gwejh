package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	nlpv1 "github.com/azkifairuz/my-skripsi-gwejh/nlp-service/internal/generated/proto/nlp/v1"
	grpcserver "github.com/azkifairuz/my-skripsi-gwejh/nlp-service/internal/grpc"
	"github.com/azkifairuz/my-skripsi-gwejh/nlp-service/internal/llm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50054"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	llmClient := llm.NewClient(llm.Config{
		BaseURL:    os.Getenv("OLLAMA_BASE_URL"),
		Model:      os.Getenv("NLP_LLM_MODEL"),
		Timeout:    envDurationSeconds("NLP_LLM_TIMEOUT_SECONDS", 10*time.Second),
		MaxRetries: envInt("NLP_LLM_MAX_RETRIES", 2),
	})

	server := grpc.NewServer()
	nlpv1.RegisterNlpServiceServer(server, grpcserver.NewServer(llmClient))
	reflection.Register(server)

	log.Printf("nlp gRPC server listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func envDurationSeconds(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}
