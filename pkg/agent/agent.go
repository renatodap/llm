package agent

import (
    i "github.com/pradord/llm/internal/agent"
    i_llm "github.com/pradord/llm/internal/llm"
    p_llm "github.com/pradord/llm/pkg/llm"
)

type (
    Executor = i.Executor
)

func NewExecutor(client *p_llm.Client) *Executor {
    return i.NewExecutor((*i_llm.Client)(client))
}

func NewExecutorWithConfig(client *p_llm.Client, maxSteps int, maxChars int, temperature float64) *Executor {
    return i.NewExecutorWithConfig((*i_llm.Client)(client), maxSteps, maxChars, temperature)
}
