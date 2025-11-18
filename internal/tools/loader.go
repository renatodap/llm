package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "github.com/pradord/llm/internal/llm"
    "gopkg.in/yaml.v3"
)

// ToolMetadata represents tool info loaded from YAML
type ToolMetadata struct {
    Name          string                 `yaml:"name"`
    Description   string                 `yaml:"description"`
    RequiredModel llm.Model              `yaml:"required_model"`
    ModelType     string                 `yaml:"model_type"`
    Parameters    map[string]interface{} `yaml:"parameters"`
}

// WrappedTool overlays metadata on top of an existing Tool, delegating Execute to the base
type WrappedTool struct {
    base        Tool
    name        string
    description string
    parameters  interface{}
    model       llm.Model
    modelType   llm.ModelType
}

func (w *WrappedTool) Name() string                 { return w.name }
func (w *WrappedTool) Description() string          { return w.description }
func (w *WrappedTool) Parameters() interface{}      { return w.parameters }
func (w *WrappedTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    return w.base.Execute(ctx, args)
}
func (w *WrappedTool) RequiredModel() llm.Model     { if w.model != "" { return w.model }; return w.base.RequiredModel() }
func (w *WrappedTool) ModelType() llm.ModelType     { if w.modelType != llm.ModelTypeInvalid { return w.modelType }; return w.base.ModelType() }

// LoadToolMetadataDir loads all *.yaml from dir and returns metadata
func LoadToolMetadataDir(dir string) ([]ToolMetadata, error) {
    var out []ToolMetadata
    if dir == "" { return out, nil }
    if _, err := os.Stat(dir); err != nil { return out, nil }
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        if !strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") { return nil }
        data, err := os.ReadFile(path)
        if err != nil { return err }
        var tm ToolMetadata
        if err := yaml.Unmarshal(data, &tm); err != nil { return fmt.Errorf("parse tool yaml %s: %w", path, err) }
        if tm.Name == "" { return fmt.Errorf("tool yaml %s missing name", path) }
        out = append(out, tm)
        return nil
    })
    return out, err
}

// ApplyToolMetadata wraps existing tools in registry with metadata if names match
func ApplyToolMetadata(reg *ToolRegistry, metas []ToolMetadata) error {
    for _, m := range metas {
        base, ok := reg.Get(m.Name)
        if !ok {
            // metadata exists but no handler; warn and continue
            continue
        }
        mt := llm.ModelTypeInvalid
        switch strings.ToLower(m.ModelType) {
        case "text": mt = llm.ModelTypeText
        case "image": mt = llm.ModelTypeImage
        case "audio": mt = llm.ModelTypeAudio
        case "video": mt = llm.ModelTypeVideo
        case "transcribe": mt = llm.ModelTypeTranscribe
        case "embedding": mt = llm.ModelTypeEmbedding
        case "vision": mt = llm.ModelTypeVision
        }
        wrapped := &WrappedTool{
            base:        base,
            name:        m.Name,
            description: m.Description,
            parameters:  m.Parameters,
            model:       m.RequiredModel,
            modelType:   mt,
        }
        reg.Register(wrapped)
    }
    return nil
}
