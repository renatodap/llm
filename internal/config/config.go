package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pradord/llm/internal/llm"
	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
    LLM   LLMConfig   `json:"llm" yaml:"llm"`
    Tools ToolsConfig `json:"tools" yaml:"tools"`
    Agent AgentConfig `json:"agent" yaml:"agent"`
    Capabilities CapabilityConfig `json:"capabilities" yaml:"capabilities"`
    Auth  AuthConfig  `json:"auth" yaml:"auth"`
    Persistence PersistenceConfig `json:"persistence" yaml:"persistence"`
    UseRealLLM bool `json:"use_real_llm" yaml:"use_real_llm"` // Toggle between mock and real OpenRouter calls
}

// LLMConfig holds LLM client configuration
type LLMConfig struct {
	APIKey          string    `json:"api_key" yaml:"api_key"`
	BaseURL         string    `json:"base_url" yaml:"base_url"`
	DefaultModel    llm.Model `json:"default_model" yaml:"default_model"`
	DefaultTemp     float64   `json:"default_temperature" yaml:"default_temperature"`
	TimeoutSeconds  int       `json:"timeout_seconds" yaml:"timeout_seconds"`
	MaxRetries      int       `json:"max_retries" yaml:"max_retries"`
	RequestsPerMin  int       `json:"requests_per_min" yaml:"requests_per_min"`
}

// ToolsConfig holds tool-specific configuration
type ToolsConfig struct {
    SearchAPIKey string `json:"search_api_key" yaml:"search_api_key"`
    ImageAPIKey  string `json:"image_api_key" yaml:"image_api_key"`
}

// AgentConfig controls the agent loop behavior
type AgentConfig struct {
    MaxSteps    int     `json:"max_steps" yaml:"max_steps"`
    MaxChars    int     `json:"max_chars" yaml:"max_chars"`
    Temperature float64 `json:"temperature" yaml:"temperature"`
}

// CapabilityConfig allows overriding capability-based model lists
type CapabilityConfig struct {
    TTSModels   []llm.Model `json:"tts_models" yaml:"tts_models"`
    VideoModels []llm.Model `json:"video_models" yaml:"video_models"`
    ImageModels []llm.Model `json:"image_models" yaml:"image_models"`
}

// Auth configuration (Supabase JWT/JWKS)
type AuthConfig struct {
    Enabled bool   `json:"enabled" yaml:"enabled"`
    JWKSURL string `json:"jwks_url" yaml:"jwks_url"`
}

// Persistence toggles to switch adapters
type PersistenceConfig struct {
    Repo    string `json:"repo" yaml:"repo"`         // file|supabase (file default)
    Vectors string `json:"vectors" yaml:"vectors"`   // memory|supabase (memory default)
    SupabaseURL string `json:"supabase_url" yaml:"supabase_url"`
    SupabaseKey string `json:"supabase_key" yaml:"supabase_key"`
    VectorTable string `json:"vector_table" yaml:"vector_table"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
    return &Config{
        LLM: LLMConfig{
            BaseURL:        "https://openrouter.ai/api/v1",
            DefaultModel:   llm.ModelGPT4oMini,
            DefaultTemp:    0.7,
            TimeoutSeconds: 60,
            MaxRetries:     3,
            RequestsPerMin: 60,
        },
        Tools: ToolsConfig{},
        Agent: AgentConfig{
            MaxSteps:    6,
            MaxChars:    24000,
            Temperature: 0.2,
        },
        Capabilities: CapabilityConfig{},
        Auth: AuthConfig{Enabled: false, JWKSURL: ""},
        Persistence: PersistenceConfig{Repo: "file", Vectors: "memory", VectorTable: "embeddings"},
        UseRealLLM: false, // Default to mock responses
    }
}

// Load loads configuration from a file
// Supports .yaml, .yml, and .json formats
func Load(path string) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// If no path provided, try to find config file
	if path == "" {
		path = findConfigFile()
		if path == "" {
			// No config file found, use defaults + env vars
			return cfg.ApplyEnv(), nil
		}
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// Parse based on extension
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse yaml config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse json config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format: %s (use .yaml, .yml, or .json)", ext)
	}

	// Apply environment variable overrides
	return cfg.ApplyEnv(), nil
}

// ApplyEnv applies environment variable overrides
func (c *Config) ApplyEnv() *Config {
	// LLM config from env
	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
		c.LLM.APIKey = apiKey
	}
	if baseURL := os.Getenv("LLM_BASE_URL"); baseURL != "" {
		c.LLM.BaseURL = baseURL
	}
	if model := os.Getenv("LLM_DEFAULT_MODEL"); model != "" {
		c.LLM.DefaultModel = llm.Model(model)
	}

	// Tool config from env
	if searchKey := os.Getenv("SEARCH_API_KEY"); searchKey != "" {
		c.Tools.SearchAPIKey = searchKey
	}
    if imageKey := os.Getenv("IMAGE_API_KEY"); imageKey != "" {
        c.Tools.ImageAPIKey = imageKey
    }

    // Toggle real LLM vs mock
    if useReal := os.Getenv("USE_REAL_LLM"); useReal == "true" || useReal == "1" {
        c.UseRealLLM = true
    }

    return c
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.LLM.APIKey == "" {
		return fmt.Errorf("LLM API key is required (set OPENROUTER_API_KEY or api_key in config)")
	}
	if c.LLM.BaseURL == "" {
		return fmt.Errorf("LLM base URL cannot be empty")
	}
	if !c.LLM.DefaultModel.IsValid() {
		return fmt.Errorf("invalid default model: %s", c.LLM.DefaultModel)
	}
	if c.LLM.DefaultTemp < 0 || c.LLM.DefaultTemp > 2 {
		return fmt.Errorf("temperature must be between 0 and 2, got: %f", c.LLM.DefaultTemp)
	}
	if c.LLM.TimeoutSeconds < 1 {
		return fmt.Errorf("timeout must be at least 1 second")
	}
	return nil
}

// findConfigFile looks for config files in common locations
func findConfigFile() string {
	// Possible config file names
	names := []string{
		"llm.yaml",
		"llm.yml",
		"llm.json",
		".llm.yaml",
		".llm.yml",
		".llm.json",
	}

	// Possible locations (in order of priority)
	locations := []string{
		".",                      // Current directory
		"config",                 // config subdirectory
		filepath.Join(os.Getenv("HOME"), ".config", "llm"), // ~/.config/llm
	}

	for _, loc := range locations {
		for _, name := range names {
			path := filepath.Join(loc, name)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// SaveExample creates an example config file
func SaveExample(path string, format string) error {
	cfg := DefaultConfig()
	cfg.LLM.APIKey = "your-api-key-here"
	cfg.Tools.SearchAPIKey = "your-search-api-key"

	var data []byte
	var err error

	switch format {
	case "yaml", "yml":
		data, err = yaml.Marshal(cfg)
	case "json":
		data, err = json.MarshalIndent(cfg, "", "  ")
	default:
		return fmt.Errorf("unsupported format: %s (use yaml or json)", format)
	}

	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}
