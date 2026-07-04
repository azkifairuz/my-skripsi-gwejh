package grpcserver

import (
	"context"

	ocrv1 "github.com/azkifairuz/my-skripsi-gwejh/ocr-service/internal/generated/proto/ocr/v1"
)

type Server struct {
	ocrv1.UnimplementedOcrServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HealthCheck(ctx context.Context, req *ocrv1.HealthCheckRequest) (*ocrv1.HealthCheckResponse, error) {
	return &ocrv1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) Ping(ctx context.Context, req *ocrv1.PingRequest) (*ocrv1.PingResponse, error) {
	return &ocrv1.PingResponse{Message: "pong"}, nil
}
