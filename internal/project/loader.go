package project

import (
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "time"
    "gopkg.in/yaml.v3"
)

// LoadDir loads all project YAMLs from a directory
func LoadDir(dir string) (map[string]*Project, error) {
    out := map[string]*Project{}
    if dir == "" { return out, nil }
    if _, err := os.Stat(dir); err != nil { return out, nil }
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        if !strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") { return nil }
        data, err := os.ReadFile(path)
        if err != nil { return err }
        var p Project
        if err := yaml.Unmarshal(data, &p); err != nil { return err }
        if p.ID == "" { p.ID = strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())) }
        now := time.Now()
        if p.CreatedAt.IsZero() { p.CreatedAt = now }
        p.UpdatedAt = now
        out[p.ID] = &p
        return nil
    })
    return out, err
}

