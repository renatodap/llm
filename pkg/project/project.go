package project

import (
    i "github.com/pradord/llm/internal/project"
)

type (
    FileStore = i.FileStore
    Project = i.Project
)

func NewFileStore(dir string) (*FileStore, error) { return i.NewFileStore(dir) }
func LoadDir(dir string) (map[string]*Project, error) { return i.LoadDir(dir) }

