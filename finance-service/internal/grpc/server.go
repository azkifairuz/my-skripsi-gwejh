package grpcserver

import (
	"context"

	financev1 "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/generated/proto/finance/v1"
)

type Server struct {
	financev1.UnimplementedFinanceServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HealthCheck(ctx context.Context, req *financev1.HealthCheckRequest) (*financev1.HealthCheckResponse, error) {
	return &financev1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) Ping(ctx context.Context, req *financev1.PingRequest) (*financev1.PingResponse, error) {
	return &financev1.PingResponse{Message: "pong"}, nil
}
