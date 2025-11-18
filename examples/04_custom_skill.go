//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pradord/llm/pkg/agent"
	"github.com/pradord/llm/pkg/llm"
	"github.com/pradord/llm/pkg/skills"
)

// Example 4: Custom Skill
// Demonstrates creating a custom skill with specialized system prompt
func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := llm.New(llm.ClientConfig{
		APIKey:       apiKey,
		DefaultModel: llm.ModelClaude35Sonnet,
	})

	// Define a custom skill
	customSkill := &skills.Skill{
		ID:          "database-architect",
		Name:        "Database Architect",
		Description: "Expert in database design and optimization",
		SystemPrompt: `You are an expert database architect with deep knowledge of:
- Relational database design (PostgreSQL, MySQL)
- NoSQL databases (MongoDB, Redis)
- Schema design and normalization
- Indexing strategies
- Query optimization
- Scalability patterns

When given a requirements description, you:
1. Analyze the data model requirements
2. Design an optimal schema
3. Recommend indexes
4. Suggest optimization strategies
5. Consider scalability from the start

Provide clear, actionable advice with code examples where appropriate.`,
		Category: "architecture",
		Tags:     []string{"database", "architecture", "optimization"},
	}

	executor := agent.NewExecutor(client)

	task := `Design a database schema for a social media platform with:
- Users (profiles, authentication)
- Posts (text, images, videos)
- Comments (nested replies)
- Likes and reactions
- Followers/following
- Direct messages

Consider: Scale to 1M users, fast feeds, efficient queries.`

	fmt.Println("Task: Design database schema for social media platform")
	fmt.Println("\n--- Database Architect working ---\n")

	result, err := executor.Run(context.Background(), customSkill, task, nil)
	if err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	fmt.Println("Result:\n", result)
}

/*
USAGE:

1. Set your API key:
   export OPENROUTER_API_KEY=your_key_here

2. Run:
   go run examples/04_custom_skill.go

OUTPUT:
Task: Design database schema for social media platform

--- Database Architect working ---

Result:
Based on your requirements, here's an optimal database schema:

1. USERS TABLE
CREATE TABLE users (
  id UUID PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  ...
);

[Detailed schema design with explanations]

2. INDEXING STRATEGY
- B-tree indexes on foreign keys
- GIN indexes for full-text search
...

3. SCALABILITY CONSIDERATIONS
- Partition posts table by date
- Read replicas for feed queries
...

CREATING YOUR OWN SKILLS:

1. Define system prompt with expertise area
2. Specify analysis methodology
3. Include formatting preferences
4. Add relevant examples
5. Test with various inputs

BUILT-IN SKILLS:
- Researcher, Writer, CodeReviewer
- TechnicalArchitect, BugAnalyzer
- TestGenerator, DatabaseDesigner
*/
//go:build examples
