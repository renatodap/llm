package agent

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"

    "github.com/pradord/llm/internal/llm"
    "github.com/pradord/llm/internal/skills"
    "github.com/pradord/llm/internal/tools"
)

// Message represents a turn in the agent conversation
// We keep a light wrapper but operate primarily on OpenAI-compatible messages

// Executor coordinates a multi-step tool-using loop with an LLM
type Executor struct {
    client       *llm.Client
    maxSteps     int
    temperature  float64
    maxChars     int // soft budget on total characters across messages
}

func NewExecutor(client *llm.Client) *Executor {
    return &Executor{client: client, maxSteps: 6, temperature: 0.2, maxChars: 24000}
}

func NewExecutorWithConfig(client *llm.Client, maxSteps int, maxChars int, temperature float64) *Executor {
    e := NewExecutor(client)
    if maxSteps > 0 { e.maxSteps = maxSteps }
    if maxChars > 0 { e.maxChars = maxChars }
    if temperature > 0 { e.temperature = temperature }
    return e
}

// Run executes a skill with iterative tool-calling
// Protocol: model can respond with final text, or JSON {"tool":"name","args":{...}}
func (e *Executor) Run(ctx context.Context, skill *skills.Skill, userPrompt string, toolList []tools.Tool) (string, error) {
    // System context with tool descriptions
    var b strings.Builder
    b.WriteString(skill.SystemPrompt)
    b.WriteString("\n\nYou can use tools to complete the task. If you choose to use a tool, respond ONLY with JSON in the form {\"tool\":\"name\",\"args\":{...}}. Otherwise, reply with the final answer in plain text.\n")
    b.WriteString("\nAvailable tools:\n")
    for _, t := range toolList {
        paramsJSON, _ := json.Marshal(t.Parameters())
        fmt.Fprintf(&b, "- %s: %s\n  parameters: %s\n", t.Name(), t.Description(), string(paramsJSON))
    }

    // Build message array per OpenAI schema
    messages := []map[string]interface{}{
        {"role": "system", "content": b.String()},
        {"role": "user", "content": userPrompt},
    }

    for step := 0; step < e.maxSteps; step++ {
        // Enforce soft character budget
        if e.totalChars(messages) > e.maxChars {
            return "Context budget exceeded before completion.", nil
        }
        // Build tool schemas
        var toolSchemas []llm.ToolFunction
        for _, t := range toolList {
            toolSchemas = append(toolSchemas, llm.ToolFunction{
                Type: "function",
                Function: map[string]interface{}{
                    "name":        t.Name(),
                    "description": t.Description(),
                    "parameters":  t.Parameters(),
                },
            })
        }
        content, calls, err := e.client.ChatWithTools(ctx, messages, toolSchemas, llm.WithModel(skill.DefaultModel), llm.WithTemperature(e.temperature))
        if err != nil { return "", err }

        if len(calls) == 0 {
            // Final answer
            if content == "" { return "", fmt.Errorf("model returned empty response") }
            messages = append(messages, map[string]interface{}{"role": "assistant", "content": content})
            return content, nil
        }
        // Append assistant message with tool_calls (content may be empty)
        // Convert toolCalls to schema array
        var toolCallsAny []map[string]interface{}
        for _, c := range calls {
            toolCallsAny = append(toolCallsAny, map[string]interface{}{
                "id":   c.ID,
                "type": "function",
                "function": map[string]interface{}{
                    "name":      c.Function.Name,
                    "arguments": c.Function.Arguments,
                },
            })
        }
        messages = append(messages, map[string]interface{}{
            "role":       "assistant",
            "content":    content,
            "tool_calls": toolCallsAny,
        })

        // Execute tool calls in order
        for _, c := range calls {
            var selected tools.Tool
            for _, t := range toolList {
                if t.Name() == c.Function.Name { selected = t; break }
            }
            if selected == nil {
                messages = append(messages, map[string]interface{}{"role": "tool", "tool_call_id": c.ID, "content": fmt.Sprintf("tool not found: %s", c.Function.Name)})
                continue
            }
            // Parse args JSON from string
            var argMap map[string]interface{}
            _ = json.Unmarshal([]byte(c.Function.Arguments), &argMap)
            argsBuf, _ := json.Marshal(argMap)
            toolOut, err := selected.Execute(ctx, argsBuf)
            if err != nil { toolOut = fmt.Sprintf("tool error: %v", err) }
            // Append tool result with tool_call_id per spec
            messages = append(messages, map[string]interface{}{
                "role":         "tool",
                "tool_call_id": c.ID,
                "content":      toolOut,
            })
            time.Sleep(100 * time.Millisecond)
        }
        // Loop for next step
    }
    return "Max steps reached without final answer.", nil
}

func (e *Executor) totalChars(msgs []map[string]interface{}) int {
    n := 0
    for _, m := range msgs {
        if s, ok := m["content"].(string); ok { n += len(s) }
    }
    return n
}

// Chat performs a simple chat completion without tools
func (e *Executor) Chat(ctx context.Context, messages []map[string]interface{}, model llm.Model) (string, error) {
    content, _, err := e.client.ChatWithTools(ctx, messages, nil, llm.WithModel(model), llm.WithTemperature(e.temperature))
    return content, err
}

// GetClient returns the underlying LLM client
func (e *Executor) GetClient() *llm.Client {
    return e.client
}
