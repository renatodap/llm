package tools

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/pradord/llm/internal/llm"
)

type WebSearch struct {
    client *llm.Client
    model  llm.Model
}

// NewWebSearch uses an LLM (via OpenRouter or Groq) to perform search-style queries.
// For OpenRouter, a search-optimized model like Perplexity Sonar can browse server-side.
func NewWebSearch(client *llm.Client, model llm.Model) *WebSearch {
    return &WebSearch{client: client, model: model}
}

func (ws *WebSearch) Name() string {
	return "web_search"
}

func (ws *WebSearch) Description() string {
    return "Search the web for current information"
}

func (ws *WebSearch) RequiredModel() llm.Model {
    return ws.model
}

func (ws *WebSearch) ModelType() llm.ModelType {
    return llm.ModelTypeText
}

func (ws *WebSearch) Parameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Search query",
			},
			"num_results": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results to return",
				"default":     5,
			},
		},
		"required": []string{"query"},
	}
}

func (ws *WebSearch) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query      string `json:"query"`
        NumResults int    `json:"num_results"`
    }

	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}

	if params.NumResults == 0 {
		params.NumResults = 5
	}

    // Delegate to LLM with a search-optimized system prompt
    prompt := fmt.Sprintf("You are a web research assistant. Search the web and summarize the top %d results about: %s. Cite sources with links.", params.NumResults, params.Query)
    out, err := ws.client.LLM(ctx, prompt, llm.WithModel(ws.model), llm.WithTemperature(0.1))
    if err != nil {
        return "", err
    }
    return out, nil
}

// SearchResult represents a single search result
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// search performs the actual web search
func (ws *WebSearch) search(ctx context.Context, query string, numResults int) []SearchResult { return nil }

// formatSearchResults converts results to a string
func formatSearchResults(results []SearchResult) string {
	if len(results) == 0 {
		return "No results found"
	}

	output := "Search Results:\n\n"
	for i, result := range results {
		output += fmt.Sprintf("%d. %s\n   URL: %s\n   %s\n\n",
			i+1, result.Title, result.URL, result.Snippet)
	}
	return output
}
