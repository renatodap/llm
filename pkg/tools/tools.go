package tools

import (
    i "github.com/pradord/llm/internal/tools"
    i_llm "github.com/pradord/llm/internal/llm"
    p_llm "github.com/pradord/llm/pkg/llm"
)

type (
    Tool = i.Tool
    ToolRegistry = i.ToolRegistry
)

func NewToolRegistry() *ToolRegistry { return i.NewToolRegistry() }

// Re-export helpers to register built-in tools through pkg API when needed
func NewWebSearch(client *p_llm.Client, model p_llm.Model) Tool { return i.NewWebSearch((*i_llm.Client)(client), model) }
