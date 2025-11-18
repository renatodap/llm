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

// AudioTTS generates speech audio from text via OpenRouter (OpenAI-compatible TTS endpoint)
type AudioTTS struct{
    apiKey  string
    baseURL string
    client  *http.Client
}

func NewAudioTTS(apiKey string) *AudioTTS { return &AudioTTS{apiKey: apiKey, baseURL: "https://openrouter.ai/api/v1", client: &http.Client{Timeout: 30 * time.Second}} }

func (t *AudioTTS) Name() string { return "generate_audio" }
func (t *AudioTTS) Description() string { return "Generate speech audio from text (returns base64 audio if supported by provider)" }
func (t *AudioTTS) RequiredModel() llm.Model { return llm.ModelElevenLabs }
func (t *AudioTTS) ModelType() llm.ModelType { return llm.ModelTypeAudio }

func (t *AudioTTS) Parameters() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "text": map[string]interface{}{"type": "string", "description": "Text to synthesize"},
            "voice": map[string]interface{}{"type": "string", "description": "Voice preset (optional)"},
            "model": map[string]interface{}{"type": "string", "description": "TTS model (provider-specific)", "default": llm.ModelElevenLabs},
            "format": map[string]interface{}{"type": "string", "description": "Audio format", "default": "mp3"},
        },
        "required": []string{"text"},
    }
}

func (t *AudioTTS) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    if t.apiKey == "" { return "", fmt.Errorf("audio generation requires OPENROUTER_API_KEY") }
    var p struct { Text string `json:"text"`; Voice string `json:"voice"`; Model string `json:"model"`; Format string `json:"format"` }
    if err := json.Unmarshal(args, &p); err != nil { return "", err }
    if p.Model == "" {
        m := llm.PickTTSCapable()
        if m == "" { return "", fmt.Errorf("no TTS-capable model configured") }
        p.Model = string(m)
    }
    if p.Format == "" { p.Format = "mp3" }

    // OpenAI-compatible TTS endpoint proxied by OpenRouter (if supported)
    body := map[string]interface{}{
        "model": p.Model,
        "input": p.Text,
        "voice": p.Voice,
        "format": p.Format,
        "response_format": "b64_json",
    }
    buf, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, "POST", t.baseURL+"/audio/speech", bytes.NewBuffer(buf))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+t.apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := t.client.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("HTTP %d from TTS API (model may not support TTS)", resp.StatusCode)
    }
    var out struct { B64 string `json:"b64_json"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if out.B64 == "" { return "", fmt.Errorf("no audio data returned") }
    // Return a data URL for convenience
    return fmt.Sprintf("data:audio/%s;base64,%s", p.Format, out.B64), nil
}
