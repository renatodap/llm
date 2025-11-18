package tools

import (
    "context"
    "encoding/json"
)

// ImageGenerator is a built-in tool
type ImageGenerator struct {
    apiKey string
}

func NewImageGenerator(apiKey string) *ImageGenerator {
    return &ImageGenerator{apiKey: apiKey}
}

func (ig *ImageGenerator) Name() string {
    return "generate_image"
}

func (ig *ImageGenerator) Description() string {
    return "Generate an image from a text description"
}

func (ig *ImageGenerator) Parameters() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "prompt": map[string]interface{}{
                "type":        "string",
                "description": "Description of the image to generate",
            },
            "size": map[string]interface{}{
                "type":        "string",
                "description": "Image size: '1024x1024', '512x512', etc",
                "default":     "1024x1024",
            },
        },
        "required": []string{"prompt"},
    }
}

func (ig *ImageGenerator) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Prompt string `json:"prompt"`
        Size   string `json:"size"`
    }
    
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    // Call DALL-E or similar
    imageURL := ig.callImageAPI(ctx, params.Prompt, params.Size)
    
    return fmt.Sprintf("Generated image: %s", imageURL), nil
}
