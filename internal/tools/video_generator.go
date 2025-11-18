package tools

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/pradord/llm/internal/llm"
)

// VideoGenerator attempts video generation via OpenRouter if supported by model/provider
type VideoGenerator struct{ apiKey, baseURL string; client *http.Client }

func NewVideoGenerator(apiKey string) *VideoGenerator {
    return &VideoGenerator{apiKey: apiKey, baseURL: "https://openrouter.ai/api/v1", client: &http.Client{Timeout: 60 * time.Second}}
}
func (vg *VideoGenerator) Name() string { return "generate_video" }
func (vg *VideoGenerator) Description() string { return "Generate a short video from a prompt (if provider supports video)" }
func (vg *VideoGenerator) RequiredModel() llm.Model { return "" }
func (vg *VideoGenerator) ModelType() llm.ModelType { return llm.ModelTypeVideo }

func (vg *VideoGenerator) Parameters() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "prompt": map[string]interface{}{"type": "string", "description": "Video prompt/story"},
            "model":  map[string]interface{}{"type": "string", "description": "Video-capable model (provider-specific)"},
            "duration": map[string]interface{}{"type": "integer", "description": "Duration seconds", "default": 5},
    },
        "required": []string{"prompt"},
    }
}

func (vg *VideoGenerator) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    if vg.apiKey == "" { return "", fmt.Errorf("video generation requires OPENROUTER_API_KEY") }
    var p struct { Prompt string `json:"prompt"`; Model string `json:"model"`; Duration int `json:"duration"` }
    if err := json.Unmarshal(args, &p); err != nil { return "", err }
    if p.Model == "" {
        m := llm.PickVideoCapable()
        if m == "" { return "", fmt.Errorf("no video-capable model available via OpenRouter") }
        p.Model = string(m)
    }

    // Hypothetical OpenAI-compatible video endpoint via OpenRouter (will return clear error if unsupported)
    body := map[string]interface{}{
        "model": p.Model,
        "prompt": p.Prompt,
        "duration": p.Duration,
        "response_format": "url",
    }
    buf, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, "POST", vg.baseURL+"/video/generations", bytes.NewBuffer(buf))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+vg.apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := vg.client.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("HTTP %d from video API (provider/model may not support video)", resp.StatusCode)
    }
    var out struct { Data []struct{ URL string `json:"url"` } `json:"data"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if len(out.Data) == 0 || out.Data[0].URL == "" { return "", fmt.Errorf("no video URL returned") }
    return out.Data[0].URL, nil
}
