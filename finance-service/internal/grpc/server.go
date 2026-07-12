package grpcserver

import (
	"context"
	"errors"
	"strings"
	"time"

	financev1 "github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/generated/proto/finance/v1"
	"github.com/azkifairuz/my-skripsi-gwejh/finance-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	financev1.UnimplementedFinanceServiceServer
	accounts     *repository.AccountRepository
	categories   *repository.CategoryRepository
	transactions *repository.TransactionRepository
}

func NewServer(accounts *repository.AccountRepository, categories *repository.CategoryRepository, transactions *repository.TransactionRepository) *Server {
	return &Server{
		accounts:     accounts,
		categories:   categories,
		transactions: transactions,
	}
}

func (s *Server) HealthCheck(ctx context.Context, req *financev1.HealthCheckRequest) (*financev1.HealthCheckResponse, error) {
	return &financev1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *Server) Ping(ctx context.Context, req *financev1.PingRequest) (*financev1.PingResponse, error) {
	return &financev1.PingResponse{Message: "pong"}, nil
}

func (s *Server) RegisterByWhatsappNumber(ctx context.Context, req *financev1.RegisterByWhatsappNumberRequest) (*financev1.RegisterByWhatsappNumberResponse, error) {
	if s.accounts == nil {
		return nil, status.Error(codes.FailedPrecondition, "database is not configured")
	}

	whatsappNumber := normalizeWhatsappNumber(req.GetWhatsappNumber())
	if whatsappNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "whatsapp_number is required")
	}

	account, err := s.accounts.RegisterByWhatsappNumber(ctx, repository.RegisterAccountParams{
		WhatsappNumber: whatsappNumber,
		Username:       strings.TrimSpace(req.GetUsername()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register account: %v", err)
	}

	return &financev1.RegisterByWhatsappNumberResponse{
		AccountId:       account.AccountID,
		WhatsappNumber:  account.WhatsappNumber,
		Username:        account.Username,
		PrimaryWalletId: account.PrimaryWalletID,
		Created:         account.Created,
		CreatedAt:       account.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) LoginByWhatsappNumber(ctx context.Context, req *financev1.LoginByWhatsappNumberRequest) (*financev1.LoginByWhatsappNumberResponse, error) {
	if s.accounts == nil {
		return nil, status.Error(codes.FailedPrecondition, "database is not configured")
	}

	whatsappNumber := normalizeWhatsappNumber(req.GetWhatsappNumber())
	if whatsappNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "whatsapp_number is required")
	}

	account, err := s.accounts.LoginByWhatsappNumber(ctx, whatsappNumber)
	if errors.Is(err, repository.ErrAccountNotFound) {
		return nil, status.Error(codes.NotFound, "account is not registered")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login account: %v", err)
	}

	return &financev1.LoginByWhatsappNumberResponse{
		AccountId:       account.AccountID,
		WhatsappNumber:  account.WhatsappNumber,
		Username:        account.Username,
		PrimaryWalletId: account.PrimaryWalletID,
		CreatedAt:       account.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) ResolveCategoryByName(ctx context.Context, req *financev1.ResolveCategoryByNameRequest) (*financev1.ResolveCategoryByNameResponse, error) {
	if s.categories == nil {
		return nil, status.Error(codes.FailedPrecondition, "database is not configured")
	}

	category, err := s.categories.ResolveByName(ctx, req.GetCategoryName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolve category: %v", err)
	}

	return &financev1.ResolveCategoryByNameResponse{
		CategoryId: category.CategoryID,
		Name:       category.Name,
	}, nil
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

func (s *Server) ListTransactionHistory(ctx context.Context, req *financev1.ListTransactionHistoryRequest) (*financev1.ListTransactionHistoryResponse, error) {
	if s.transactions == nil {
		return nil, status.Error(codes.FailedPrecondition, "database is not configured")
	}

	params, err := validateListTransactionHistoryRequest(req)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactions.ListHistory(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list transaction history: %v", err)
	}

	items := make([]*financev1.TransactionHistoryItem, 0, len(transactions))
	for _, transaction := range transactions {
		items = append(items, &financev1.TransactionHistoryItem{
			TransactionId:   transaction.TransactionID,
			AccountId:       transaction.AccountID,
			WalletId:        transaction.WalletID,
			WalletName:      transaction.WalletName,
			CategoryId:      transaction.CategoryID,
			CategoryName:    transaction.CategoryName,
			Type:            transaction.Type,
			Amount:          transaction.Amount,
			Name:            transaction.Name,
			IsAiGenerated:   transaction.IsAIGenerated,
			ReceiptImageUrl: transaction.ReceiptImageURL,
			ReportDate:      transaction.ReportDate.Format(time.RFC3339),
			CreatedAt:       transaction.CreatedAt.Format(time.RFC3339),
		})
	}

	return &financev1.ListTransactionHistoryResponse{
		Transactions: items,
		Limit:        params.Limit,
		Offset:       params.Offset,
		Count:        int32(len(items)),
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

func validateListTransactionHistoryRequest(req *financev1.ListTransactionHistoryRequest) (repository.ListTransactionHistoryParams, error) {
	if strings.TrimSpace(req.GetAccountId()) == "" {
		return repository.ListTransactionHistoryParams{}, status.Error(codes.InvalidArgument, "account_id is required")
	}

	transactionType := strings.TrimSpace(req.GetType())
	if transactionType != "" && transactionType != "income" && transactionType != "expense" {
		return repository.ListTransactionHistoryParams{}, status.Error(codes.InvalidArgument, "type must be income, expense, or empty")
	}

	limit := req.GetLimit()
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := req.GetOffset()
	if offset < 0 {
		return repository.ListTransactionHistoryParams{}, status.Error(codes.InvalidArgument, "offset must be greater than or equal to 0")
	}

	fromDate, err := parseOptionalRFC3339(req.GetFromDate(), "from_date")
	if err != nil {
		return repository.ListTransactionHistoryParams{}, err
	}

	toDate, err := parseOptionalRFC3339(req.GetToDate(), "to_date")
	if err != nil {
		return repository.ListTransactionHistoryParams{}, err
	}

	return repository.ListTransactionHistoryParams{
		AccountID: strings.TrimSpace(req.GetAccountId()),
		Limit:     limit,
		Offset:    offset,
		Type:      transactionType,
		FromDate:  fromDate,
		ToDate:    toDate,
	}, nil
}

func parseOptionalRFC3339(value, field string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s must use RFC3339 format", field)
	}

	return &parsed, nil
}

func normalizeWhatsappNumber(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "+")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	return value
}
