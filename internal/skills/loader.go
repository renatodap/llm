package skills

import (
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "github.com/pradord/llm/internal/llm"
    "github.com/pradord/llm/internal/tools"
    "gopkg.in/yaml.v3"
)

// SkillYAML captures skill metadata for YAML loading
type SkillYAML struct {
    Name         string     `yaml:"name"`
    Description  string     `yaml:"description"`
    SystemPrompt string     `yaml:"system_prompt"`
    Tools        []string   `yaml:"tools"`
    Resources    []string   `yaml:"resources"`
    Examples     []string   `yaml:"examples"`
    DefaultModel llm.Model  `yaml:"default_model"`
}

// LoadSkillsDir loads skills from YAML files and registers them
func LoadSkillsDir(dir string, reg *SkillRegistry) error {
    if dir == "" { return nil }
    if _, err := os.Stat(dir); err != nil { return nil }
    return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        if !strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") { return nil }
        data, err := os.ReadFile(path)
        if err != nil { return err }
        var s SkillYAML
        if err := yaml.Unmarshal(data, &s); err != nil { return err }
        if s.Name == "" { return nil }
        reg.Register(&Skill{
            Name:         s.Name,
            Description:  s.Description,
            SystemPrompt: s.SystemPrompt,
            Tools:        s.Tools,
            Resources:    s.Resources,
            Examples:     s.Examples,
            DefaultModel: s.DefaultModel,
        })
        return nil
    })
}

// Optionally ensure tools exist for a skill
func ValidateSkillTools(reg *SkillRegistry, toolReg *tools.ToolRegistry) map[string][]string {
    missing := make(map[string][]string)
    for _, sk := range reg.List() {
        for _, tn := range sk.Tools {
            if _, ok := toolReg.Get(tn); !ok {
                missing[sk.Name] = append(missing[sk.Name], tn)
            }
        }
    }
    return missing
}

