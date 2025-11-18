package embeddings

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Client calls an OpenAI-compatible embeddings endpoint (e.g., via OpenRouter)
type Client struct {
    apiKey  string
    baseURL string
    http    *http.Client
    model   string
}

type Config struct {
    APIKey  string
    BaseURL string
    Model   string // e.g., openai/text-embedding-ada-002
    Timeout time.Duration
}

func New(cfg Config) *Client {
    if cfg.BaseURL == "" { cfg.BaseURL = "https://openrouter.ai/api/v1" }
    if cfg.Timeout == 0 { cfg.Timeout = 30 * time.Second }
    if cfg.Model == "" { cfg.Model = "openai/text-embedding-ada-002" }
    return &Client{apiKey: cfg.APIKey, baseURL: cfg.BaseURL, model: cfg.Model, http: &http.Client{Timeout: cfg.Timeout}}
}

// Embed returns an embedding vector for the given input text
func (c *Client) Embed(ctx context.Context, input string) ([]float64, error) {
    body := map[string]interface{}{
        "model": c.model,
        "input": input,
    }
    buf, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/embeddings", bytes.NewBuffer(buf))
    if err != nil { return nil, err }
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.http.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("unexpected status %d from embeddings API", resp.StatusCode)
    }
    var out struct{ Data []struct{ Embedding []float64 `json:"embedding"` } `json:"data"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return nil, err }
    if len(out.Data) == 0 { return nil, fmt.Errorf("no embedding returned") }
    return out.Data[0].Embedding, nil
}

// Cosine similarity
func Cosine(a, b []float64) float64 {
    if len(a) == 0 || len(b) == 0 || len(a) != len(b) { return 0 }
    var dot, na, nb float64
    for i := range a { dot += a[i]*b[i]; na += a[i]*a[i]; nb += b[i]*b[i] }
    if na == 0 || nb == 0 { return 0 }
    return dot / (sqrt(na) * sqrt(nb))
}

func sqrt(x float64) float64 {
    // minimal sqrt to avoid extra imports
    z := x
    if z == 0 { return 0 }
    for i := 0; i < 20; i++ { z = 0.5*(z + x/z) }
    return z
}

