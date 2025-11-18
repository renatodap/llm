package tools

// GetSafeDefaultTools returns a list of safe tools that can be auto-registered
// These tools don't have dangerous side effects (no file writes, code execution, etc.)
func GetSafeDefaultTools() []Tool {
    return []Tool{
        // WebSearch now requires an LLM client; register manually in main
        NewURLFetcher(),       // Safe - just reads URLs
        NewCalculator(),       // Safe - just math
        NewImageAnalyzer(""), // Requires API key from user
    }
}

// GetAllBuiltInTools returns all available built-in tools
// Some of these may require explicit user opt-in for safety
func GetAllBuiltInTools() map[string]Tool {
    return map[string]Tool{
        // Safe tools
        // "web_search" requires an LLM client; register manually in main
        "fetch_url":     NewURLFetcher(),
        "calculator":    NewCalculator(),
        "analyze_image": NewImageAnalyzer(""),
        "generate_image": NewImageGenerator("", ""),
        "generate_audio": NewAudioTTS(""),
        "generate_video": NewVideoGenerator(""),

        // Potentially dangerous - require opt-in
        "execute_code": NewCodeExecutor(),
    }
}

// AutoRegisterSafeTools registers safe default tools to a registry
// This is called automatically unless the user provides their own tools
func AutoRegisterSafeTools(registry *ToolRegistry, _searchAPIKey, imageAPIKey string) {
	// WebSearch requires an LLM client; register in main where client exists
	registry.Register(NewURLFetcher())
	registry.Register(NewCalculator())
	registry.Register(NewImageAnalyzer(imageAPIKey))
}
