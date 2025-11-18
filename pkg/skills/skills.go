package skills

import (
    i "github.com/pradord/llm/internal/skills"
    i_tools "github.com/pradord/llm/internal/tools"
    p_tools "github.com/pradord/llm/pkg/tools"
)

type (
    Skill = i.Skill
    SkillRegistry = i.SkillRegistry
)

func NewSkillRegistry(toolRegistry *p_tools.ToolRegistry) *SkillRegistry {
    return i.NewSkillRegistry((*i_tools.ToolRegistry)(toolRegistry))
}

// Simple re-exports for common built-in skills
func NewResearchAssistant() *Skill { return i.NewResearchAssistant() }
func NewContentCreator() *Skill { return i.NewContentCreator() }
func NewCodeReviewer() *Skill { return i.NewCodeReviewer() }
func NewDataAnalyst() *Skill { return i.NewDataAnalyst() }
func NewCoder() *Skill { return i.NewCoder() }
func NewSummarizer() *Skill { return i.NewSummarizer() }
func NewTranslator() *Skill { return i.NewTranslator() }
func NewTutor() *Skill { return i.NewTutor() }
