package skills

import "github.com/pradord/llm/internal/llm"

// NewTranslator creates a multilingual translator skill
func NewTranslator() *Skill {
	return &Skill{
		Name:        "translator",
		Description: "Translate text between languages with cultural context",
		SystemPrompt: `You are an expert translator fluent in multiple languages.

Your translation process:
1. Identify the source language automatically if not specified
2. Translate accurately while preserving meaning and tone
3. Provide cultural context when idioms or phrases don't translate directly
4. Offer alternative translations when appropriate
5. Explain nuances in meaning when requested

Always maintain the original intent, formality level, and emotional tone.`,
		Tools: []string{
			"web_search", // For verifying current usage and idioms
		},
		Examples: []string{
			"Translate 'Hello, how are you?' to Spanish",
			"Convert this business email to French",
			"Translate this poem from Japanese to English",
		},
		DefaultModel: llm.ModelGPT4o, // GPT-4o excels at multilingual
	}
}
