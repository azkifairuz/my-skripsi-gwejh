package client

import (
	"context"

	nlpv1 "github.com/azkifairuz/my-skripsi-gwejh/whatsapp-service/internal/generated/proto/nlp/v1"
)

func (d *Dependencies) PingNLP(ctx context.Context) (string, error) {
	conn, err := dial(d.nlpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := nlpv1.NewNlpServiceClient(conn).Ping(ctx, &nlpv1.PingRequest{
		Source: requestSource,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
