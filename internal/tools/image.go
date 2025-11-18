package tools

import (
    "bytes"
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/pradord/llm/internal/llm"
)

// ImageAnalyzer analyzes images using vision models
type ImageAnalyzer struct {
    apiKey  string
    client  *http.Client
    baseURL string
}

// NewImageAnalyzer creates a new image analyzer
func NewImageAnalyzer(apiKey string) *ImageAnalyzer {
    return &ImageAnalyzer{apiKey: apiKey, client: &http.Client{Timeout: 20 * time.Second}, baseURL: "https://openrouter.ai/api/v1"}
}

func (ia *ImageAnalyzer) Name() string {
	return "analyze_image"
}

func (ia *ImageAnalyzer) Description() string {
	return "Analyze an image and describe its contents"
}

func (ia *ImageAnalyzer) RequiredModel() llm.Model {
	return llm.ModelGPT4o // Needs vision capabilities
}

func (ia *ImageAnalyzer) ModelType() llm.ModelType {
	return llm.ModelTypeText // Uses vision model but returns text
}

func (ia *ImageAnalyzer) Parameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"image_url": map[string]interface{}{
				"type":        "string",
				"description": "URL or base64-encoded image",
			},
			"prompt": map[string]interface{}{
				"type":        "string",
				"description": "Question or instruction about the image",
				"default":     "Describe this image in detail",
			},
		},
		"required": []string{"image_url"},
	}
}

func (ia *ImageAnalyzer) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		ImageURL string `json:"image_url"`
		Prompt   string `json:"prompt"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}

	if params.Prompt == "" {
		params.Prompt = "Describe this image in detail"
	}

	// Analyze the image
	result := ia.analyzeImage(ctx, params.ImageURL, params.Prompt)

	return result, nil
}

// analyzeImage performs the actual image analysis
func (ia *ImageAnalyzer) analyzeImage(ctx context.Context, imageURL, prompt string) string {
    if ia.apiKey == "" {
        return "Image analysis requires IMAGE_API_KEY or OPENROUTER_API_KEY."
    }
    // Prepare OpenRouter request for multimodal content (text + image_url)
    body := map[string]interface{}{
        "model": llm.ModelGPT4o.String(),
        "messages": []map[string]interface{}{
            {
                "role": "user",
                "content": []map[string]interface{}{
                    {"type": "text", "text": prompt},
                    {"type": "image_url", "image_url": map[string]interface{}{"url": imageURL}},
                },
            },
        },
        "temperature": 0.2,
    }
    buf, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, "POST", ia.baseURL+"/chat/completions", bytes.NewBuffer(buf))
    if err != nil {
        return fmt.Sprintf("request error: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+ia.apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := ia.client.Do(req)
    if err != nil {
        return fmt.Sprintf("http error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Sprintf("HTTP %d during vision analysis", resp.StatusCode)
    }
    var result struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Sprintf("decode error: %v", err)
    }
    if len(result.Choices) == 0 {
        return "no response from model"
    }
    return result.Choices[0].Message.Content
}

// isURL checks if a string is a URL
func isURL(s string) bool {
	return len(s) > 7 && (s[:7] == "http://" || s[:8] == "https://")
}

// encodeImage encodes an image to base64
func encodeImage(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}
