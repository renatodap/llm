package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pradord/llm/internal/llm"
)

// URLFetcher fetches content from URLs
type URLFetcher struct {
	httpClient *http.Client
}

func NewURLFetcher() *URLFetcher {
	return &URLFetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (uf *URLFetcher) Name() string {
	return "fetch_url"
}

func (uf *URLFetcher) Description() string {
	return "Fetch and read content from a URL"
}

func (uf *URLFetcher) Parameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "URL to fetch content from",
			},
		},
		"required": []string{"url"},
	}
}

func (uf *URLFetcher) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "LLM-Library/1.0")

	resp, err := uf.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Limit response size to 1MB
	limitedReader := io.LimitReader(resp.Body, 1024*1024)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	return fmt.Sprintf("Content from %s:\n\n%s", params.URL, string(body)), nil
}

func (uf *URLFetcher) RequiredModel() llm.Model {
	return "" // No model needed
}

func (uf *URLFetcher) ModelType() llm.ModelType {
	return llm.ModelTypeText
}
