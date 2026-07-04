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
