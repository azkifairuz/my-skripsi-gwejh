package client

import (
	"context"

	ocrv1 "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/generated/proto/ocr/v1"
)

func (d *Dependencies) PingOCR(ctx context.Context) (string, error) {
	conn, err := dial(d.ocrAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := ocrv1.NewOcrServiceClient(conn).Ping(ctx, &ocrv1.PingRequest{
		Source: requestSource,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
