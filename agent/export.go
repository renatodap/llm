// Package agent provides a thin public wrapper re-exporting the internal agent
// types for orchestrating multi-step tool-using executions.
package agent

import i "github.com/pradord/llm/internal/agent"
import l "github.com/pradord/llm/internal/llm"

// Re-exports
type (
    Executor = i.Executor
)

// Constructors
func NewExecutor(client *l.Client) *Executor { return i.NewExecutor(client) }
func NewExecutorWithConfig(client *l.Client, maxSteps, maxChars int, temperature float64) *Executor {
    return i.NewExecutorWithConfig(client, maxSteps, maxChars, temperature)
}
