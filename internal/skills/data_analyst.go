package skills

import "github.com/pradord/llm/internal/llm"

// NewDataAnalyst creates a data analysis skill
func NewDataAnalyst() *Skill {
	return &Skill{
		Name:        "data_analyst",
		Description: "Analyze data, find patterns, and provide insights",
		SystemPrompt: `You are a data analyst skilled at extracting insights from data.

Your analysis process:
1. Understand the data structure and context
2. Identify key patterns and trends
3. Perform statistical analysis
4. Create visualizations (describe them clearly)
5. Provide actionable insights and recommendations

Always explain your methodology and reasoning clearly.`,
		Tools: []string{
			// Data analysis tools would go here
		},
		Examples: []string{
			"Analyze this sales data for trends",
			"Find correlations in this user behavior dataset",
			"Provide insights from this customer feedback",
		},
		DefaultModel: llm.ModelGPT4o,
	}
}
