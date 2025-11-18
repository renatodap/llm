package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/pradord/llm/pkg/config"
    "github.com/pradord/llm/pkg/conversation"
    "github.com/pradord/llm/pkg/llm"
    "github.com/pradord/llm/pkg/skills"
    "github.com/pradord/llm/pkg/tools"
    "github.com/pradord/llm/pkg/agent"
    "github.com/pradord/llm/pkg/project"
)

func main() {
    // Command line flags
    configPath := flag.String("config", "", "Path to config file (auto-detected if not provided)")
    generateConfig := flag.String("generate-config", "", "Generate example config file (yaml or json)")
    // Agent overrides and defs dir
    agentSteps := flag.Int("agent-steps", 0, "Override agent max steps")
    agentMaxChars := flag.Int("agent-max-chars", 0, "Override agent max chars budget")
    agentTemp := flag.Float64("agent-temp", 0, "Override agent temperature")
    defsDir := flag.String("defs-dir", "", "Path to defs directory containing tools/ and skills/ YAML")
    projectsDir := flag.String("projects-dir", "", "Path to projects directory with YAML definitions")
    projectID := flag.String("project", "", "Select a project ID to scope tools/skills and system prompt")
    // Retrieval flags
    threadsDir := flag.String("threads-dir", ".llm_threads", "Directory to store conversation threads")
    threadID := flag.String("thread-id", "", "Conversation thread ID for retrieval context")
    retrieve := flag.Bool("retrieve", true, "Enable retrieval of prior messages as context")
    flag.Parse()

    // Load projects if provided and pick selected project
    var activeProject *project.Project
    if *projectsDir != "" {
        projs, err := project.LoadDir(*projectsDir)
        if err != nil {
            fmt.Printf("Warning: project defs load error: %v\n", err)
        } else if *projectID != "" {
            activeProject = projs[*projectID]
        }
    }

    // Generate example config if requested
    if *generateConfig != "" {
        format := "yaml"
        if *generateConfig == "json" {
            format = "json"
        }
        filename := "llm.example." + format
        if err := config.SaveExample(filename, format); err != nil {
            fmt.Printf("Error generating config: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Generated example config: %s\n", filename)
        fmt.Println("Copy this file to llm.yaml and customize it")
        return
    }

    fmt.Println("=== LLM Library with Tools & Skills Demo ===\n")

    // Load configuration
    cfg, err := config.Load(*configPath)
    if err != nil {
        fmt.Printf("Error loading config: %v\n", err)
        os.Exit(1)
    }

    // Validate configuration
    if err := cfg.Validate(); err != nil {
        fmt.Printf("Invalid config: %v\n", err)
        fmt.Println("\nTip: Set OPENROUTER_API_KEY environment variable or create llm.yaml config file")
        fmt.Println("Run with --generate-config=yaml to create an example config")
        os.Exit(1)
    }

    fmt.Printf("Config loaded (using model: %s)\n\n", cfg.LLM.DefaultModel)

    // Create LLM client from config
    client := llm.NewClient(llm.ClientConfig{
        APIKey:         cfg.LLM.APIKey,
        BaseURL:        cfg.LLM.BaseURL,
        DefaultModel:   cfg.LLM.DefaultModel,
        DefaultTemp:    cfg.LLM.DefaultTemp,
        TimeoutSeconds: cfg.LLM.TimeoutSeconds,
        MaxRetries:     cfg.LLM.MaxRetries,
        RequestsPerMin: cfg.LLM.RequestsPerMin,
    })

    // Optionally override capability lists from config
    if len(cfg.Capabilities.TTSModels) > 0 || len(cfg.Capabilities.VideoModels) > 0 || len(cfg.Capabilities.ImageModels) > 0 {
        llm.UpdateCapabilityModels(cfg.Capabilities.TTSModels, cfg.Capabilities.VideoModels, cfg.Capabilities.ImageModels)
    }

    // Create tool registry
    toolRegistry := tools.NewToolRegistry()

    // Register tools with config
    toolRegistry.Register(tools.NewWebSearch(client, llm.ModelPerplexitySonar))
    imageAPIKey := cfg.Tools.ImageAPIKey
    if imageAPIKey == "" {
        imageAPIKey = cfg.LLM.APIKey // Use main API key if not specified
    }
    toolRegistry.Register(tools.NewImageAnalyzer(imageAPIKey))
    // Also register safe tools used by some skills
    toolRegistry.Register(tools.NewURLFetcher())
    toolRegistry.Register(tools.NewCalculator())
    // Register generation tools
    toolRegistry.Register(tools.NewImageGenerator(imageAPIKey, cfg.LLM.BaseURL))
    toolRegistry.Register(tools.NewAudioTTS(cfg.Tools.ImageAPIKey))
    toolRegistry.Register(tools.NewVideoGenerator(cfg.Tools.ImageAPIKey))

    // Create skill registry
    skillRegistry := skills.NewSkillRegistry(toolRegistry)

    // Register built-in skills
    skillRegistry.Register(skills.NewResearchAssistant())
    skillRegistry.Register(skills.NewContentCreator())
    skillRegistry.Register(skills.NewCodeReviewer())
    skillRegistry.Register(skills.NewDataAnalyst())

    // Load YAML tool/skill definitions if provided
    if *defsDir != "" {
        // Tools metadata: defs/tools/*.yaml
        if metas, err := tools.LoadToolMetadataDir(filepath.Join(*defsDir, "tools")); err == nil {
            // Apply metadata to existing handlers (wrap)
            _ = tools.ApplyToolMetadata(toolRegistry, metas)
        } else {
            fmt.Printf("Warning: tool defs load error: %v\n", err)
        }
        // Skills: defs/skills/*.yaml
        if err := skills.LoadSkillsDir(filepath.Join(*defsDir, "skills"), skillRegistry); err != nil {
            fmt.Printf("Warning: skill defs load error: %v\n", err)
        }
        // Optional: validate tool references
        if missing := skills.ValidateSkillTools(skillRegistry, toolRegistry); len(missing) > 0 {
            fmt.Printf("Warning: some skills reference missing tools: %v\n", missing)
        }
    }

    // Example 1: Using client with config defaults
    fmt.Println("=== Example 1: Simple LLM Call (using config defaults) ===")
    response, err := client.LLM(
        context.Background(),
        "Explain Go in one sentence",
        // No WithModel needed - uses default from config
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Response: %s\n", response)
    }
    fmt.Println()

    // Example 2: Override config defaults
    fmt.Println("=== Example 2: Override Config Defaults ===")
    response, err = client.LLM(
        context.Background(),
        "What is the capital of France?",
        llm.WithModel(llm.ModelClaude35Sonnet), // Override default model
        llm.WithTemperature(0.3),                // Override default temperature
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Response: %s\n", response)
    }
    fmt.Println()

    // Example 3: Show available tools with model requirements
    fmt.Println("=== Example 3: Available Tools ===")
    toolSchemas := toolRegistry.List()
    for _, schema := range toolSchemas {
        funcData := schema["function"].(map[string]interface{})
        toolName := funcData["name"].(string)
        tool, _ := toolRegistry.Get(toolName)

        fmt.Printf("Tool: %s\n", funcData["name"])
        fmt.Printf("Description: %s\n", funcData["description"])
        if reqModel := tool.RequiredModel(); reqModel != "" {
            fmt.Printf("Required Model: %s\n", reqModel)
        }
        fmt.Printf("Model Type: %s\n", tool.ModelType())
        fmt.Println()
    }

    // Example 4: Show available skills
    fmt.Println("=== Example 4: Available Skills ===")
    for _, skill := range skillRegistry.List() {
        fmt.Printf("Skill: %s\n", skill.Name)
        fmt.Printf("Description: %s\n", skill.Description)
        fmt.Printf("Default Model: %s\n", skill.DefaultModel)
        fmt.Printf("Uses Tools: %v\n", skill.Tools)
        fmt.Println()
    }

    // Example 5: Execute a tool directly
    fmt.Println("=== Example 5: Execute Web Search Tool ===")
    searchArgs, _ := json.Marshal(map[string]interface{}{
        "query":       "golang best practices 2025",
        "num_results": 3,
    })
    result, err := toolRegistry.Execute(context.Background(), "web_search", searchArgs)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println(result)
    }

    // Example 6: Show configuration in use
    fmt.Println("=== Example 6: Current Configuration ===")
    fmt.Printf("Base URL: %s\n", cfg.LLM.BaseURL)
    fmt.Printf("Default Model: %s\n", cfg.LLM.DefaultModel)
    fmt.Printf("Default Temperature: %.1f\n", cfg.LLM.DefaultTemp)
    fmt.Printf("Timeout: %d seconds\n", cfg.LLM.TimeoutSeconds)
    fmt.Printf("Max Retries: %d\n", cfg.LLM.MaxRetries)
    fmt.Printf("Rate Limit: %d requests/min\n", cfg.LLM.RequestsPerMin)
    fmt.Println()

    // Example 7: Skill Execution via agent loop (multi-step)
    fmt.Println("=== Example 7: Execute Skill via Agent Loop ===")
    // Build executor with config/flags
    maxSteps := cfg.Agent.MaxSteps
    maxChars := cfg.Agent.MaxChars
    temp := cfg.Agent.Temperature
    if *agentSteps > 0 { maxSteps = *agentSteps }
    if *agentMaxChars > 0 { maxChars = *agentMaxChars }
    if *agentTemp > 0 { temp = *agentTemp }
    exec := agent.NewExecutorWithConfig(client, maxSteps, maxChars, temp)
    // Choose skill and apply project system prompt if any
    skill := skills.NewResearchAssistant()
    if activeProject != nil && activeProject.SystemPrompt != "" {
        skill.SystemPrompt = activeProject.SystemPrompt + "\n\n" + skill.SystemPrompt
    }
    // Simple retrieval: use last N messages as context (no vector store needed)
    if *retrieve && *threadID != "" {
        fs, err := conversation.NewFileStore(*threadsDir)
        if err == nil {
            if th, err := fs.GetThread(*threadID); err == nil {
                // Use last 5 messages as context
                start := len(th.Messages) - 5
                if start < 0 {
                    start = 0
                }
                if len(th.Messages) > 0 {
                    var ctxBlock string
                    for i := start; i < len(th.Messages); i++ {
                        m := th.Messages[i]
                        ctxBlock += fmt.Sprintf("- (%s) %s\n", m.Role, m.Content)
                    }
                    if ctxBlock != "" {
                        skill.SystemPrompt = skill.SystemPrompt + "\n\nRecent Context:\n" + ctxBlock
                    }
                }
            }
        }
    }
    // Scope tools if project specifies
    toolNames := []string{"web_search"}
    if activeProject != nil && len(activeProject.Tools) > 0 { toolNames = activeProject.Tools }
    userPrompt := "Summarize the latest stable features of Go language in 3 bullets."
    out, err := exec.Run(context.Background(), skill, userPrompt, mustGetTools(toolRegistry, toolNames))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println(out)
    }

    // Append messages to thread (if thread provided)
    if *threadID != "" {
        if fs, err := conversation.NewFileStore(*threadsDir); err == nil {
            th, err := fs.GetThread(*threadID)
            if err == nil {
                _, _ = fs.AppendMessage(th.ID, conversation.Message{Role: "user", Content: userPrompt})
                _, _ = fs.AppendMessage(th.ID, conversation.Message{Role: "assistant", Content: out})
            }
        }
    }

    fmt.Println("=== Demo Complete ===")
    fmt.Println("\nConfiguration Priority:")
    fmt.Println("1. Environment variables (highest)")
    fmt.Println("2. Config file (llm.yaml or llm.json)")
    fmt.Println("3. Built-in defaults (lowest)")
    fmt.Println("\nTo customize: Create llm.yaml in current directory")
    fmt.Println("Run: ./llm --generate-config=yaml")
}

// helper to resolve tool implementations from names
func mustGetTools(reg *tools.ToolRegistry, names []string) []tools.Tool {
    out := make([]tools.Tool, 0, len(names))
    for _, n := range names {
        if t, ok := reg.Get(n); ok {
            out = append(out, t)
        }
    }
    return out
}
