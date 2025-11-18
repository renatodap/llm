package skills

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pradord/llm/internal/llm"
	"github.com/pradord/llm/internal/tools"
)

// Skill orchestrates tools with specific instructions and domain knowledge
type Skill struct {
	Name           string
	Description    string
	SystemPrompt   string
	Tools          []string  // Tool names this skill can use
	Resources      []string  // Files, docs, context
	Examples       []string  // Few-shot examples
	RequiredSkills []string  // Other skills this skill depends on
	DefaultModel   llm.Model
}

// SkillRegistry manages available skills
type SkillRegistry struct {
	skills map[string]*Skill
	tools  *tools.ToolRegistry
}

// NewSkillRegistry creates a new skill registry
func NewSkillRegistry(toolRegistry *tools.ToolRegistry) *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]*Skill),
		tools:  toolRegistry,
	}
}

// Register adds a skill to the registry
func (sr *SkillRegistry) Register(skill *Skill) {
	sr.skills[skill.Name] = skill
}

// Get retrieves a skill by name
func (sr *SkillRegistry) Get(name string) (*Skill, bool) {
	skill, ok := sr.skills[name]
	return skill, ok
}

// List returns all registered skills
func (sr *SkillRegistry) List() []*Skill {
	skillList := make([]*Skill, 0, len(sr.skills))
	for _, skill := range sr.skills {
		skillList = append(skillList, skill)
	}
	return skillList
}

// GetTools returns all tools that a skill uses
func (sr *SkillRegistry) GetTools(skillName string) ([]tools.Tool, error) {
	skill, ok := sr.Get(skillName)
	if !ok {
		return nil, fmt.Errorf("skill not found: %s", skillName)
	}

	var toolList []tools.Tool
	for _, toolName := range skill.Tools {
		tool, ok := sr.tools.Get(toolName)
		if !ok {
			return nil, fmt.Errorf("tool not found for skill %s: %s", skillName, toolName)
		}
		toolList = append(toolList, tool)
	}
	return toolList, nil
}

// Execute runs a skill with the given user prompt
func (sr *SkillRegistry) Execute(ctx context.Context, skillName, userPrompt string, executor SkillExecutor) (string, error) {
	skill, ok := sr.Get(skillName)
	if !ok {
		return "", fmt.Errorf("skill not found: %s", skillName)
	}

	// Get the tools this skill needs
	toolList, err := sr.GetTools(skillName)
	if err != nil {
		return "", err
	}

	// Execute the skill with its context
	return executor.ExecuteSkill(ctx, skill, userPrompt, toolList)
}

// SkillExecutor is the interface that LLM clients must implement to execute skills
type SkillExecutor interface {
	ExecuteSkill(ctx context.Context, skill *Skill, userPrompt string, tools []tools.Tool) (string, error)
}

// SkillResponse represents the result of executing a skill
type SkillResponse struct {
	Result       string
	ToolsUsed    []string
	TokensUsed   int
	ModelUsed    llm.Model
	ResourcesRef []string
}

// ToJSON converts the response to JSON
func (sr *SkillResponse) ToJSON() (string, error) {
	data, err := json.MarshalIndent(sr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
