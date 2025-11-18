// Package skills re-exports the internal skills types and constructors so that
// applications can define and register skills without importing internal paths.
package skills

import i "github.com/pradord/llm/internal/skills"
import t "github.com/pradord/llm/internal/tools"

// Re-export core types
type (
    Skill         = i.Skill
    SkillRegistry = i.SkillRegistry
)

// Constructors
func NewSkillRegistry(toolRegistry *t.ToolRegistry) *SkillRegistry { return i.NewSkillRegistry(toolRegistry) }

// Built-in skills
func NewResearchAssistant() *Skill { return i.NewResearchAssistant() }
func NewContentCreator() *Skill { return i.NewContentCreator() }
func NewCodeReviewer() *Skill { return i.NewCodeReviewer() }
func NewDataAnalyst() *Skill { return i.NewDataAnalyst() }
func NewTranslator() *Skill { return i.NewTranslator() }
func NewSummarizer() *Skill { return i.NewSummarizer() }
func NewTutor() *Skill { return i.NewTutor() }
func NewCoder() *Skill { return i.NewCoder() }

// All methods remain available on the aliased type.
