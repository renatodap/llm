# LLM Library Examples

Complete working examples demonstrating different use cases of the LLM library.

## Prerequisites

1. **Get an OpenRouter API key**: https://openrouter.ai/keys
2. **Set environment variable**:
   ```bash
   export OPENROUTER_API_KEY=your_key_here
   ```

## Examples

### 01. Simple Chat
**File**: `01_simple_chat.go`

Basic LLM usage with a single prompt. Perfect for getting started.

```bash
go run examples/01_simple_chat.go
```

**What you'll learn**:
- Creating an LLM client
- Making basic API calls
- Using different models

---

### 02. Streaming Response
**File**: `02_streaming_response.go`

Real-time streaming of LLM responses for chat interfaces.

```bash
go run examples/02_streaming_response.go
```

**What you'll learn**:
- Streaming responses chunk-by-chunk
- Implementing real-time UIs
- Managing stream lifecycle

---

### 03. Agent with Tools
**File**: `03_agent_with_tools.go`

Multi-step reasoning with tool calling (web search, calculator, etc.).

```bash
go run examples/03_agent_with_tools.go
```

**What you'll learn**:
- Creating agent executors
- Defining tools
- Multi-step reasoning
- Using skills

---

### 04. Custom Skill
**File**: `04_custom_skill.go`

Creating specialized skills with custom system prompts.

```bash
go run examples/04_custom_skill.go
```

**What you'll learn**:
- Defining custom skills
- Writing effective system prompts
- Domain expertise injection

---

### 05. Batch Processing
**File**: `05_batch_processing.go`

Processing multiple prompts efficiently in parallel.

```bash
go run examples/05_batch_processing.go
```

**What you'll learn**:
- Parallel processing
- Goroutines with LLM calls
- Batch optimization
- Cost-effective models

---

### 06. (Removed) REST API Client
Server/API functionality has been removed in this repository. Use this repo as a pure library and expose your own HTTP endpoints in a separate application if needed.

---

## Mock Mode (No API Key Needed)

All examples work in mock mode for testing:

```bash
USE_REAL_LLM=false go run examples/01_simple_chat.go
```

Mock mode returns predefined responses without calling the API.

## Running All Examples

```bash
# Run each example
for example in examples/*.go; do
    echo "Running $example..."
    go run "$example"
done
```

## Integration with Your Project

### As a Library

```go
import "github.com/pradord/llm/pkg/llm"

client := llm.New(llm.ClientConfig{
    APIKey: os.Getenv("OPENROUTER_API_KEY"),
    DefaultModel: llm.ModelClaude35Sonnet,
})

response, _ := client.LLM(ctx, "Your prompt")
```

### As a REST API
Not included. Build your own server that imports this library.

## Next Steps

1. **Try the examples** - Run each one to understand the capabilities
2. **Read the docs** - See [../docs/SETUP.md](../docs/SETUP.md) for setup
3. **Build something** - Use as a library or copy components into your project

## Model Selection Guide

### For Development (Free)
- `ModelLlama3170B` - Fast, free via Groq
- `ModelMixtral` - Good quality, free

### For Production (Paid)
- `ModelClaude35Sonnet` - Best quality, $3/$15 per 1M tokens
- `ModelGPT4o` - Fast, reliable, $5/$15 per 1M tokens

### For Batch Processing
- Use free models (Llama, Mixtral)
- Implement rate limiting
- Monitor costs

## Troubleshooting

### "API key is required" Error
```bash
# Make sure to export, not just set
export OPENROUTER_API_KEY=your_key_here

# Verify it's set
echo $OPENROUTER_API_KEY
```

### "Connection refused" Error
Make sure the server is running:
```bash
./llm_server.exe --addr=:3001
```

### Rate Limiting
OpenRouter has rate limits. For high-volume:
- Use free models
- Implement backoff
- Consider batch endpoints

## Support

- API Documentation: [../docs/API.md](../docs/API.md)
- OpenRouter Docs: https://openrouter.ai/docs
- Main README: [../README.md](../README.md)
