package conversation

import (
    i "github.com/pradord/llm/internal/conversation"
)

type (
    FileStore = i.FileStore
    Message = i.Message
    Thread = i.Thread
)

func NewFileStore(dir string) (*FileStore, error) { return i.NewFileStore(dir) }

