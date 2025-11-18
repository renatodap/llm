package config

import (
    i "github.com/pradord/llm/internal/config"
    p_llm "github.com/pradord/llm/pkg/llm"
)

type (
    Config = i.Config
    LLMConfig = i.LLMConfig
    ToolsConfig = i.ToolsConfig
    AgentConfig = i.AgentConfig
    CapabilityConfig = i.CapabilityConfig
    AuthConfig = i.AuthConfig
    PersistenceConfig = i.PersistenceConfig
)

// Re-map model type fields for LLMConfig to pkg llm.Model through type aliasing
// (Already handled by aliasing in internal definitions)

func DefaultConfig() *Config { return i.DefaultConfig() }
func Load(path string) (*Config, error) { return i.Load(path) }
func SaveExample(path string, format string) error { return i.SaveExample(path, format) }

// Convenience: ensure DefaultModel type compiles with pkg llm.Model
func (c *Config) SetDefaultModel(m p_llm.Model) { c.LLM.DefaultModel = m }

