package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	requestSource = "whatsapp-service"
	callTimeout   = 3 * time.Second
)

type Dependencies struct {
	financeAddr string
	nlpAddr     string
	ocrAddr     string
}

type PingResult struct {
	Finance string
	NLP     string
	OCR     string
}

func NewDependencies(financeAddr, nlpAddr, ocrAddr string) *Dependencies {
	return &Dependencies{
		financeAddr: financeAddr,
		nlpAddr:     nlpAddr,
		ocrAddr:     ocrAddr,
	}
}

func (d *Dependencies) PingAll(ctx context.Context) (*PingResult, error) {
	finance, err := d.PingFinance(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping finance service: %w", err)
	}

	nlp, err := d.PingNLP(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping nlp service: %w", err)
	}

	ocr, err := d.PingOCR(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping ocr service: %w", err)
	}

	return &PingResult{
		Finance: finance,
		NLP:     nlp,
		OCR:     ocr,
	}, nil
}

func dial(target string) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, callTimeout)
}
