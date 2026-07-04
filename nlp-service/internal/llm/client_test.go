package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestExtractTransaction(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://ollama.test",
		Model:   "qwen:4b",
		Timeout: time.Second,
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/api/generate" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}

			var body bytes.Buffer
			if err := json.NewEncoder(&body).Encode(generateResponse{
				Response: `{"type":"expense","amount":20000,"category":"makanan","description":"bayar kopi"}`,
			}); err != nil {
				t.Fatalf("encode response: %v", err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(&body),
			}, nil
		}),
	}

	transaction, err := client.ExtractTransaction(context.Background(), "bayar kopi 20 ribu")
	if err != nil {
		t.Fatalf("ExtractTransaction() error = %v", err)
	}

	if transaction.Type != "expense" {
		t.Fatalf("Type = %q, want expense", transaction.Type)
	}
	if transaction.Amount != 20000 {
		t.Fatalf("Amount = %d, want 20000", transaction.Amount)
	}
	if transaction.Category != "makanan" {
		t.Fatalf("Category = %q, want makanan", transaction.Category)
	}
	if transaction.Description != "bayar kopi" {
		t.Fatalf("Description = %q, want bayar kopi", transaction.Description)
	}
}

func TestExtractTransactionWithoutAmount(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://example.test",
		Timeout: time.Second,
	})

	_, err := client.ExtractTransaction(context.Background(), "bayar kopi")
	if err == nil {
		t.Fatal("ExtractTransaction() error = nil, want error")
	}
}

func TestExtractTransactionFallsBackWhenModelReturnsPlaceholder(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://ollama.test",
		Model:   "qwen:4b",
		Timeout: time.Second,
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			var body bytes.Buffer
			if err := json.NewEncoder(&body).Encode(generateResponse{
				Response: `{"type":"income|expense","amount":0,"category":"makanan|transportasi","description":"transaction description without money amount"}`,
			}); err != nil {
				t.Fatalf("encode response: %v", err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(&body),
			}, nil
		}),
	}

	transaction, err := client.ExtractTransaction(context.Background(), "bayar kopi 20 ribu")
	if err != nil {
		t.Fatalf("ExtractTransaction() error = %v", err)
	}

	if transaction.Type != "expense" || transaction.Amount != 20000 || transaction.Category != "makanan" || transaction.Description != "bayar kopi" {
		t.Fatalf("ExtractTransaction() = %+v", transaction)
	}
}

func TestFallbackExtractTransactionWithJuta(t *testing.T) {
	transaction, err := fallbackExtractTransaction("gaji freelance 1,5 juta")
	if err != nil {
		t.Fatalf("fallbackExtractTransaction() error = %v", err)
	}

	if transaction.Type != "income" || transaction.Amount != 1500000 || transaction.Category != "pemasukan" || transaction.Description != "gaji freelance" {
		t.Fatalf("fallbackExtractTransaction() = %+v", transaction)
	}
}

func TestParseTransactionWithExtraText(t *testing.T) {
	transaction, err := parseTransaction("hasil:\n{\"type\":\"income\",\"amount\":500000,\"category\":\"pemasukan\",\"description\":\"bonus\"}\nselesai")
	if err != nil {
		t.Fatalf("parseTransaction() error = %v", err)
	}

	if transaction.Type != "income" || transaction.Amount != 500000 || transaction.Category != "pemasukan" || transaction.Description != "bonus" {
		t.Fatalf("parseTransaction() = %+v", transaction)
	}
}

func TestValidateRejectsInvalidCategory(t *testing.T) {
	transaction := &Transaction{
		Type:        "expense",
		Amount:      10000,
		Category:    "misteri",
		Description: "bayar sesuatu",
	}

	if err := transaction.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
