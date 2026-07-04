package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "http://ollama:11434"
	defaultModel   = "qwen:4b"
)

var moneyPattern = regexp.MustCompile(`(?i)(rp\s*)?\d+([.,]\d+)?\s*(ribu|rb|k|juta|jt)?`)

type Client struct {
	baseURL    string
	model      string
	maxRetries int
	httpClient *http.Client
}

type Config struct {
	BaseURL    string
	Model      string
	Timeout    time.Duration
	MaxRetries int
}

type Transaction struct {
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type generateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Format  string                 `json:"format"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type generateResponse struct {
	Response string `json:"response"`
}

func NewClient(cfg Config) *Client {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = defaultModel
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		baseURL:    baseURL,
		model:      model,
		maxRetries: cfg.MaxRetries,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ExtractTransaction(ctx context.Context, text string) (*Transaction, error) {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return nil, errors.New("transaction text is empty")
	}
	if !moneyPattern.MatchString(normalized) {
		return nil, errors.New("transaction amount is not found")
	}

	transaction, err := c.extractTransactionOnce(ctx, normalized)
	if err == nil {
		return transaction, nil
	}

	var lastErr error = err
	for i := 0; i < c.maxRetries; i++ {
		transaction, err = c.extractTransactionOnce(ctx, normalized)
		if err == nil {
			return transaction, nil
		}
		lastErr = err
	}

	transaction, fallbackErr := fallbackExtractTransaction(normalized)
	if fallbackErr == nil {
		return transaction, nil
	}

	return nil, lastErr
}

func (c *Client) extractTransactionOnce(ctx context.Context, text string) (*Transaction, error) {
	reqBody, err := json.Marshal(generateRequest{
		Model:  c.model,
		Prompt: buildPrompt(text),
		Stream: false,
		Format: "json",
		Options: map[string]interface{}{
			"temperature": 0,
			"num_predict": 128,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ollama response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var generated generateResponse
	if err := json.Unmarshal(body, &generated); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	transaction, err := parseTransaction(generated.Response)
	if err != nil {
		return nil, err
	}
	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func buildPrompt(text string) string {
	return fmt.Sprintf(`You are a JSON API for Indonesian personal finance transaction extraction.
Return exactly one valid JSON object. Do not return markdown. Do not explain.
Do not copy instructions. Do not output placeholder values.
Field "type" must be exactly "expense" or "income".
Field "amount" must be an integer rupiah amount.
Field "category" must be exactly one of: makanan, transportasi, belanja, tagihan, pemasukan, kesehatan, pendidikan, hiburan, lainnya.
Field "description" must be the transaction description without money amount.

Aturan:
- Nominal harus selalu dikonversi ke integer rupiah.
- Jika teks tidak menyebut pemasukan, default type adalah expense.
- Description adalah deskripsi transaksi tanpa nominal uang.
- Jika kategori tidak jelas, gunakan lainnya.
- Untuk makanan/minuman seperti kopi, nasi, ayam, mie, bakso, gunakan kategori makanan.
- Untuk gaji, bonus, dapat uang, atau terima uang, gunakan type income dan kategori pemasukan.

Example input: bayar kopi 20 ribu
Example output: {"type":"expense","amount":20000,"category":"makanan","description":"bayar kopi"}

Example input: gaji freelance 2 juta
Example output: {"type":"income","amount":2000000,"category":"pemasukan","description":"gaji freelance"}

Input: %s
Output:`, text)
}

func fallbackExtractTransaction(text string) (*Transaction, error) {
	amountMatch := moneyPattern.FindString(text)
	if amountMatch == "" {
		return nil, errors.New("transaction amount is not found")
	}

	amount, err := parseAmount(amountMatch)
	if err != nil {
		return nil, err
	}

	description := strings.TrimSpace(moneyPattern.ReplaceAllString(text, ""))
	description = strings.Join(strings.Fields(description), " ")
	if description == "" {
		return nil, errors.New("transaction description is empty")
	}

	transaction := &Transaction{
		Type:        inferType(description),
		Amount:      amount,
		Category:    inferCategory(description),
		Description: description,
	}
	if err := transaction.Validate(); err != nil {
		return nil, err
	}
	return transaction, nil
}

func parseAmount(raw string) (int64, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "rp", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	multiplier := int64(1)
	for _, suffix := range []struct {
		text       string
		multiplier int64
	}{
		{"ribu", 1000},
		{"rb", 1000},
		{"k", 1000},
		{"juta", 1000000},
		{"jt", 1000000},
	} {
		if strings.HasSuffix(normalized, suffix.text) {
			multiplier = suffix.multiplier
			normalized = strings.TrimSuffix(normalized, suffix.text)
			break
		}
	}

	if strings.Contains(normalized, ".") && !strings.Contains(normalized, ",") {
		parts := strings.Split(normalized, ".")
		if len(parts[len(parts)-1]) == 3 {
			normalized = strings.ReplaceAll(normalized, ".", "")
		}
	}
	normalized = strings.ReplaceAll(normalized, ",", ".")

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("parse transaction amount: %w", err)
	}

	amount := int64(value * float64(multiplier))
	if amount <= 0 {
		return 0, errors.New("transaction amount must be greater than zero")
	}
	return amount, nil
}

func inferType(text string) string {
	lowered := strings.ToLower(text)
	for _, keyword := range []string{"gaji", "dapat", "terima", "masuk", "bonus", "dibayar"} {
		if strings.Contains(lowered, keyword) {
			return "income"
		}
	}
	return "expense"
}

func inferCategory(text string) string {
	lowered := strings.ToLower(text)
	categories := []struct {
		category string
		keywords []string
	}{
		{"pemasukan", []string{"gaji", "bonus", "dapat", "terima"}},
		{"makanan", []string{"makan", "ayam", "geprek", "kopi", "nasi", "minum", "bakso", "mie"}},
		{"transportasi", []string{"gojek", "grab", "bensin", "parkir", "tol", "ojek"}},
		{"belanja", []string{"beli", "shopee", "tokopedia", "barang"}},
		{"tagihan", []string{"listrik", "air", "wifi", "internet", "pulsa", "token"}},
	}

	for _, category := range categories {
		for _, keyword := range category.keywords {
			if strings.Contains(lowered, keyword) {
				return category.category
			}
		}
	}
	return "lainnya"
}

func parseTransaction(raw string) (*Transaction, error) {
	cleaned := strings.TrimSpace(raw)
	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start >= 0 && end >= start {
		cleaned = cleaned[start : end+1]
	}

	var transaction Transaction
	if err := json.Unmarshal([]byte(cleaned), &transaction); err != nil {
		return nil, fmt.Errorf("decode transaction json: %w", err)
	}
	return &transaction, nil
}

func (t *Transaction) Validate() error {
	t.Type = strings.ToLower(strings.TrimSpace(t.Type))
	t.Category = strings.ToLower(strings.TrimSpace(t.Category))
	t.Description = strings.TrimSpace(t.Description)

	if t.Type != "income" && t.Type != "expense" {
		return fmt.Errorf("invalid transaction type %q", t.Type)
	}
	if t.Amount <= 0 {
		return errors.New("transaction amount must be greater than zero")
	}
	if !validCategories[t.Category] {
		return fmt.Errorf("invalid transaction category %q", t.Category)
	}
	if t.Description == "" {
		return errors.New("transaction description is empty")
	}
	return nil
}

var validCategories = map[string]bool{
	"makanan":      true,
	"transportasi": true,
	"belanja":      true,
	"tagihan":      true,
	"pemasukan":    true,
	"kesehatan":    true,
	"pendidikan":   true,
	"hiburan":      true,
	"lainnya":      true,
}
