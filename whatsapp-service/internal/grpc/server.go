package grpcserver

import (
	"context"

	"github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/client"
	whatsappv1 "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/generated/proto/whatsapp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	whatsappv1.UnimplementedWhatsappServiceServer
	dependencies *client.Dependencies
}

func NewServer(dependencies *client.Dependencies) *Server {
	return &Server{dependencies: dependencies}
}

func (s *Server) HealthCheck(ctx context.Context, req *whatsappv1.HealthCheckRequest) (*whatsappv1.HealthCheckResponse, error) {
	return &whatsappv1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) PingDependencies(ctx context.Context, req *whatsappv1.PingDependenciesRequest) (*whatsappv1.PingDependenciesResponse, error) {
	result, err := s.dependencies.PingAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "ping dependencies failed: %v", err)
	}

	return &whatsappv1.PingDependenciesResponse{
		Finance: result.Finance,
		Nlp:     result.NLP,
		Ocr:     result.OCR,
	}, nil
}
