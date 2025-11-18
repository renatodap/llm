package skills

import "github.com/pradord/llm/internal/llm"

// NewContentCreator creates a content creator skill
func NewContentCreator() *Skill {
	return &Skill{
		Name:        "content_creator",
		Description: "Create engaging content with images and research",
		SystemPrompt: `You are a content creator who produces high-quality, engaging content.

Your workflow:
1. Research the topic thoroughly using web_search
2. Generate relevant images using analyze_image if needed
3. Write compelling copy based on research
4. Format for the target platform (blog, social, video)
5. Include visuals and cite sources

Always create content that is informative, engaging, and well-sourced.`,
		Tools: []string{
			"web_search",
			"analyze_image",
		},
		Examples: []string{
			"Write a blog post about sustainable living tips",
			"Create social media content about productivity hacks",
			"Draft a video script about the future of renewable energy",
		},
		DefaultModel: llm.ModelGPT4o,
	}
}
