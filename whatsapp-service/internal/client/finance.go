package client

import (
	"context"

	financev1 "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/generated/proto/finance/v1"
)

func (d *Dependencies) PingFinance(ctx context.Context) (string, error) {
	conn, err := dial(d.financeAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := financev1.NewFinanceServiceClient(conn).Ping(ctx, &financev1.PingRequest{
		Source: requestSource,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}

func (d *Dependencies) RegisterByWhatsappNumber(ctx context.Context, whatsappNumber, username string) (*financev1.RegisterByWhatsappNumberResponse, error) {
	conn, err := dial(d.financeAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	return financev1.NewFinanceServiceClient(conn).RegisterByWhatsappNumber(ctx, &financev1.RegisterByWhatsappNumberRequest{
		WhatsappNumber: whatsappNumber,
		Username:       username,
	})
}

func (d *Dependencies) LoginByWhatsappNumber(ctx context.Context, whatsappNumber string) (*financev1.LoginByWhatsappNumberResponse, error) {
	conn, err := dial(d.financeAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	return financev1.NewFinanceServiceClient(conn).LoginByWhatsappNumber(ctx, &financev1.LoginByWhatsappNumberRequest{
		WhatsappNumber: whatsappNumber,
	})
}

func (d *Dependencies) ListTransactionHistory(ctx context.Context, req *financev1.ListTransactionHistoryRequest) (*financev1.ListTransactionHistoryResponse, error) {
	conn, err := dial(d.financeAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	return financev1.NewFinanceServiceClient(conn).ListTransactionHistory(ctx, req)
}
