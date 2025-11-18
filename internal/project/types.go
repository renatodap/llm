package project

import "time"

// Project defines a scoped workspace with its own system prompt, tools, and skills
type Project struct {
    ID           string    `yaml:"id" json:"id"`
    Name         string    `yaml:"name" json:"name"`
    SystemPrompt string    `yaml:"system_prompt" json:"system_prompt"`
    Tools        []string  `yaml:"tools" json:"tools"`
    Skills       []string  `yaml:"skills" json:"skills"`
    DefaultModel string    `yaml:"default_model" json:"default_model"`
    OwnerUserID  string    `yaml:"owner_user_id" json:"owner_user_id"`
    CreatedAt    time.Time `yaml:"created_at" json:"created_at"`
    UpdatedAt    time.Time `yaml:"updated_at" json:"updated_at"`
}
