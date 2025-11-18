# Setup Guide - Internal LLM Toolkit

Quick setup for internal use. No database, no auth, no complications.

## Prerequisites

1. **Go 1.21+**
2. **OpenRouter API Key** (https://openrouter.ai/keys)

## Installation

### Option 1: Use as Library

```bash
# In your Go project
go get github.com/pradord/llm
```

### Option 2: Clone Repo

```bash
git clone https://github.com/pradord/llm
cd llm
```

### Option 3: Copy Components

```bash
# Just copy what you need
cp -r /path/to/llm/internal/llm your-project/
cp -r /path/to/llm/internal/agent your-project/
```

## Configuration

### 1. Environment Variables

Create `.env` file (optional):
```bash
OPENROUTER_API_KEY=your_key_here
USE_REAL_LLM=true
LLM_DEFAULT_MODEL=anthropic/claude-3.5-sonnet
```

### 2. Config File (optional)

Create `llm.yaml`:
```yaml
llm:
  apiKey: "your_openrouter_api_key"
  defaultModel: "anthropic/claude-3.5-sonnet"
  useRealLLM: true

agent:
  maxSteps: 6
  maxChars: 24000
  temperature: 0.2
```

## Running as Server

### Build

```bash
go build -o llm_server main.go
```

### Run

```bash
# Mock mode (no API key needed)
./llm_server --addr=:3001

# Real mode
USE_REAL_LLM=true OPENROUTER_API_KEY=xxx ./llm_server --addr=:3001
```

### Test

```bash
# Health check
curl http://localhost:3001/health

# Send message
curl -X POST http://localhost:3001/api/chats/test/messages \
  -H "Content-Type: application/json" \
  -d '{"content":"What is quantum computing?"}'
```

## Using as Library

### Simple Chat

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/pradord/llm/internal/llm"
)

func main() {
    client := llm.New(llm.ClientConfig{
        APIKey:       os.Getenv("OPENROUTER_API_KEY"),
        DefaultModel: llm.ModelClaude35Sonnet,
    })

    response, err := client.LLM(context.Background(), "Explain quantum computing")
    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### Agent with Skills

```go
import (
    "github.com/pradord/llm/internal/agent"
    "github.com/pradord/llm/internal/skills"
)

executor := agent.NewExecutor(client)
skill := skills.NewCodeReviewer()

result, err := executor.Run(ctx, skill, "Review this code...", nil)
```

## Examples

See [examples/](../examples/) directory:

```bash
# Run examples
cd examples

# Simple chat
OPENROUTER_API_KEY=xxx go run 01_simple_chat.go

# Streaming
OPENROUTER_API_KEY=xxx go run 02_streaming_response.go

# Agent with tools
OPENROUTER_API_KEY=xxx go run 03_agent_with_tools.go

# Custom skills
OPENROUTER_API_KEY=xxx go run 04_custom_skill.go

# Batch processing
OPENROUTER_API_KEY=xxx go run 05_batch_processing.go

# REST API client
go run 06_rest_api_client.go
```

## File Storage

Conversations are stored in `.llm_threads/`:
```
.llm_threads/
├── {uuid-1}.json
├── {uuid-2}.json
└── ...
```

No database setup needed!

## Troubleshooting

### "API key is required" Error

```bash
# Make sure to export
export OPENROUTER_API_KEY=your_key_here

# Verify
echo $OPENROUTER_API_KEY
```

### "Connection refused" Error

Make sure server is running:
```bash
./llm_server --addr=:3001
```

### Mock vs Real

```bash
# Mock mode (no API calls, returns mock responses)
USE_REAL_LLM=false ./llm_server

# Real mode (calls OpenRouter API)
USE_REAL_LLM=true OPENROUTER_API_KEY=xxx ./llm_server
```

## Development

### Build

```bash
go build -o llm_server main.go
```

### Test (coming soon)

```bash
go test ./...
```

### Add Custom Skills

```go
// internal/skills/custom.go
package skills

func NewMyCustomSkill() *Skill {
    return &Skill{
        ID:   "my-custom",
        Name: "My Custom Skill",
        SystemPrompt: `You are an expert in...`,
    }
}
```

### Add Custom Tools

```go
// internal/tools/custom.go
package tools

type MyTool struct{}

func (t *MyTool) Name() string { return "my_tool" }
func (t *MyTool) Description() string { return "Does something cool" }
func (t *MyTool) Parameters() interface{} { /* schema */ }
func (t *MyTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    // Implementation
}
```

## Production Checklist

- [ ] Get OpenRouter API key
- [ ] Set environment variables
- [ ] Test with mock mode first
- [ ] Switch to real mode
- [ ] Monitor `.llm_threads/` disk usage
- [ ] Optional: Setup logrotate or cleanup script

## Support

- API Reference: [API.md](API.md)
- Examples: [examples/](../examples/)
- OpenRouter Docs: https://openrouter.ai/docs

---

**Simple setup. No complexity.**
