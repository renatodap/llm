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

// ImageGenerator generates images from text prompts
type ImageGenerator struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

func NewImageGenerator(apiKey, baseURL string) *ImageGenerator {
    if baseURL == "" {
        baseURL = "https://openrouter.ai/api/v1"
    }
    return &ImageGenerator{apiKey: apiKey, baseURL: baseURL, client: &http.Client{Timeout: 20 * time.Second}}
}

func (ig *ImageGenerator) Name() string        { return "generate_image" }
func (ig *ImageGenerator) Description() string { return "Generate an image from a text prompt" }
func (ig *ImageGenerator) RequiredModel() llm.Model { return llm.ModelDALLE3 }
func (ig *ImageGenerator) ModelType() llm.ModelType { return llm.ModelTypeImage }

func (ig *ImageGenerator) Parameters() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "prompt": map[string]interface{}{"type": "string", "description": "Image prompt"},
            "model":  map[string]interface{}{"type": "string", "description": "Model name (optional)", "default": llm.ModelDALLE3},
            "size":   map[string]interface{}{"type": "string", "description": "Image size e.g. 1024x1024", "default": "1024x1024"},
        },
        "required": []string{"prompt"},
    }
}

func (ig *ImageGenerator) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    var p struct {
        Prompt string     `json:"prompt"`
        Model  llm.Model  `json:"model"`
        Size   string     `json:"size"`
    }
    if err := json.Unmarshal(args, &p); err != nil {
        return "", err
    }
    if p.Model == "" {
        p.Model = llm.ModelDALLE3
    }
    if ig.apiKey == "" {
        return "", fmt.Errorf("image generation requires IMAGE_API_KEY or OPENROUTER_API_KEY")
    }

    // OpenRouter-compatible image generation payload (mirrors OpenAI images)
    body := map[string]interface{}{
        "model":  p.Model,
        "prompt": p.Prompt,
        "size":   p.Size,
    }
    buf, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, "POST", ig.baseURL+"/images", bytes.NewBuffer(buf))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+ig.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := ig.client.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("HTTP %d from image API", resp.StatusCode)
    }
    var out struct{
        Data []struct{ URL string `json:"url"` } `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if len(out.Data) == 0 || out.Data[0].URL == "" {
        return "", fmt.Errorf("no image URL returned")
    }
    return fmt.Sprintf("Image URL: %s", out.Data[0].URL), nil
}

