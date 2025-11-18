//go:build examples

package main

import (
	"fmt"

    "github.com/pradord/llm/pkg/llm"
    "github.com/pradord/llm/pkg/skills"
)

// ExampleCustomSkill shows how to create a custom skill
func main() {
	// Create a custom skill
	personalTrainer := &skills.Skill{
		Name:        "personal_trainer",
		Description: "Fitness and workout planning assistant",
		SystemPrompt: `You are a certified personal trainer and nutritionist.

Your approach:
1. Assess the user's fitness level and goals
2. Create personalized workout plans
3. Provide nutrition advice
4. Track progress and adjust plans
5. Motivate and encourage

Use web_search to find exercise demonstrations and nutrition info.
Use calculator for calorie and macro calculations.`,
		Tools: []string{
			"web_search",
			"calculator",
		},
		Examples: []string{
			"Create a beginner workout plan for weight loss",
			"Calculate my daily calorie needs",
			"Suggest exercises for building muscle",
		},
		DefaultModel: llm.ModelClaude35Sonnet,
	}

	// Register it
	registry := skills.NewSkillRegistry(nil)
	registry.Register(personalTrainer)

    fmt.Println("Custom skill registered:", personalTrainer.Name)
}
