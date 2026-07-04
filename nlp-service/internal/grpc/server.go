package grpcserver

import (
	"context"

	nlpv1 "github.com/azkifairuz/my-skripsi-gwejh/nlp-service/internal/generated/proto/nlp/v1"
	"github.com/azkifairuz/my-skripsi-gwejh/nlp-service/internal/llm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	nlpv1.UnimplementedNlpServiceServer
	llmClient *llm.Client
}

func NewServer(llmClient *llm.Client) *Server {
	return &Server{
		llmClient: llmClient,
	}
}

func (s *Server) HealthCheck(ctx context.Context, req *nlpv1.HealthCheckRequest) (*nlpv1.HealthCheckResponse, error) {
	return &nlpv1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) Ping(ctx context.Context, req *nlpv1.PingRequest) (*nlpv1.PingResponse, error) {
	return &nlpv1.PingResponse{Message: "pong"}, nil
}

func (s *Server) ExtractTransaction(ctx context.Context, req *nlpv1.ExtractTransactionRequest) (*nlpv1.ExtractTransactionResponse, error) {
	if s.llmClient == nil {
		return nil, status.Error(codes.FailedPrecondition, "llm client is not configured")
	}

	transaction, err := s.llmClient.ExtractTransaction(ctx, req.GetText())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &nlpv1.ExtractTransactionResponse{
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Category:    transaction.Category,
		Description: transaction.Description,
	}, nil
}
