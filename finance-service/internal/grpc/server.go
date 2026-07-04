package grpcserver

import (
	"context"
	"strings"
	"time"

	financev1 "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/generated/proto/finance/v1"
	"github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	financev1.UnimplementedFinanceServiceServer
	transactions *repository.TransactionRepository
}

func NewServer(transactions *repository.TransactionRepository) *Server {
	return &Server{transactions: transactions}
}

func (s *Server) HealthCheck(ctx context.Context, req *financev1.HealthCheckRequest) (*financev1.HealthCheckResponse, error) {
	return &financev1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) Ping(ctx context.Context, req *financev1.PingRequest) (*financev1.PingResponse, error) {
	return &financev1.PingResponse{Message: "pong"}, nil
}

func (s *Server) CreateTransaction(ctx context.Context, req *financev1.CreateTransactionRequest) (*financev1.CreateTransactionResponse, error) {
	if s.transactions == nil {
		return nil, status.Error(codes.FailedPrecondition, "database is not configured")
	}

	reportDate, err := validateCreateTransactionRequest(req)
	if err != nil {
		return nil, err
	}

	transaction, err := s.transactions.Create(ctx, repository.CreateTransactionParams{
		AccountID:       req.GetAccountId(),
		WalletID:        req.GetWalletId(),
		CategoryID:      req.GetCategoryId(),
		BudgetID:        req.BudgetId,
		Type:            req.GetType(),
		Amount:          req.GetAmount(),
		Name:            req.GetName(),
		IsAIGenerated:   req.GetIsAiGenerated(),
		ReceiptImageURL: req.GetReceiptImageUrl(),
		ReportDate:      reportDate,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create transaction: %v", err)
	}

	return &financev1.CreateTransactionResponse{
		TransactionId: transaction.TransactionID,
		AccountId:     transaction.AccountID,
		WalletId:      transaction.WalletID,
		CategoryId:    transaction.CategoryID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Name:          transaction.Name,
		ReportDate:    transaction.ReportDate.Format(time.RFC3339),
		CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
	}, nil
}

func validateCreateTransactionRequest(req *financev1.CreateTransactionRequest) (time.Time, error) {
	if strings.TrimSpace(req.GetAccountId()) == "" {
		return time.Time{}, status.Error(codes.InvalidArgument, "account_id is required")
	}

	if req.GetWalletId() <= 0 {
		return time.Time{}, status.Error(codes.InvalidArgument, "wallet_id must be greater than 0")
	}

	if req.GetCategoryId() <= 0 {
		return time.Time{}, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}

	transactionType := req.GetType()
	if transactionType != "income" && transactionType != "expense" {
		return time.Time{}, status.Error(codes.InvalidArgument, "type must be income or expense")
	}

	if req.GetAmount() <= 0 {
		return time.Time{}, status.Error(codes.InvalidArgument, "amount must be greater than 0")
	}

	if strings.TrimSpace(req.GetName()) == "" {
		return time.Time{}, status.Error(codes.InvalidArgument, "name is required")
	}

	if strings.TrimSpace(req.GetReportDate()) == "" {
		return time.Now().UTC(), nil
	}

	reportDate, err := time.Parse(time.RFC3339, req.GetReportDate())
	if err != nil {
		return time.Time{}, status.Error(codes.InvalidArgument, "report_date must use RFC3339 format")
	}

	return reportDate, nil
}
