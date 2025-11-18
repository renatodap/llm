package conversation

import (
    "time"
)

// Message represents a single chat message
type Message struct {
    ID        string    `json:"id" yaml:"id"`
    Role      string    `json:"role" yaml:"role"` // system|user|assistant|tool
    Content   string    `json:"content" yaml:"content"`
    CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

// Thread is a chat-like conversation with ordered messages and optional summary
type Thread struct {
    ID        string                 `json:"id" yaml:"id"`
    ProjectID string                 `json:"project_id" yaml:"project_id"`
    Title     string                 `json:"title" yaml:"title"`
    CreatedAt time.Time              `json:"created_at" yaml:"created_at"`
    UpdatedAt time.Time              `json:"updated_at" yaml:"updated_at"`
    Summary   string                 `json:"summary" yaml:"summary"`
    Metadata  map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
    Messages  []Message              `json:"messages" yaml:"messages"`
}

// Store defines persistence for threads
type Store interface {
    CreateThread(title string) (*Thread, error)
    CreateThreadForProject(projectID, title string) (*Thread, error)
    GetThread(id string) (*Thread, error)
    ListThreads() ([]*Thread, error)
    AppendMessage(threadID string, msg Message) (*Thread, error)
    UpdateSummary(threadID string, summary string) error
    UpdateThread(t *Thread) error
}
