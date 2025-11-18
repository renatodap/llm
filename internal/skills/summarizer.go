package skills

import "github.com/pradord/llm/internal/llm"

// NewSummarizer creates a text summarization skill
func NewSummarizer() *Skill {
	return &Skill{
		Name:        "summarizer",
		Description: "Summarize long texts, articles, or documents concisely",
		SystemPrompt: `You are an expert at creating clear, concise summaries of complex content.

Your summarization process:
1. Read and understand the full content
2. Identify the main ideas and key points
3. Extract the most important information
4. Organize into a logical, coherent summary
5. Maintain accuracy while being concise
6. Use bullet points for clarity when appropriate

Use fetch_url to retrieve web articles when given URLs.`,
		Tools: []string{
			"fetch_url",
			"web_search",
		},
		Examples: []string{
			"Summarize this article: [URL]",
			"Give me the key points from this document",
			"Create a one-paragraph summary of this research paper",
		},
		DefaultModel: llm.ModelClaude35Sonnet,
	}
}
