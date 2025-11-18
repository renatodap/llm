package project

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

type FileStore struct {
    dir string
    mu  sync.Mutex
}

func NewFileStore(dir string) (*FileStore, error) {
    if dir == "" { dir = ".llm_projects" }
    if err := os.MkdirAll(dir, 0o755); err != nil { return nil, err }
    return &FileStore{dir: dir}, nil
}

func (fs *FileStore) path(id string) string { return filepath.Join(fs.dir, id+".json") }

func (fs *FileStore) Create(p *Project) (*Project, error) {
    fs.mu.Lock(); defer fs.mu.Unlock()
    if p.ID == "" { p.ID = strings.TrimSpace(newID()) }
    now := time.Now()
    p.CreatedAt = now; p.UpdatedAt = now
    if err := fs.save(p); err != nil { return nil, err }
    return p, nil
}

func (fs *FileStore) Get(id string) (*Project, error) {
    data, err := os.ReadFile(fs.path(id))
    if err != nil { return nil, err }
    var p Project
    if err := json.Unmarshal(data, &p); err != nil { return nil, err }
    return &p, nil
}

func (fs *FileStore) List() ([]*Project, error) {
    entries, err := os.ReadDir(fs.dir)
    if err != nil { return nil, err }
    var out []*Project
    for _, e := range entries {
        if e.IsDir() || filepath.Ext(e.Name()) != ".json" { continue }
        data, err := os.ReadFile(filepath.Join(fs.dir, e.Name()))
        if err != nil { continue }
        var p Project
        if json.Unmarshal(data, &p) == nil { out = append(out, &p) }
    }
    return out, nil
}

func (fs *FileStore) save(p *Project) error {
    data, err := json.MarshalIndent(p, "", "  ")
    if err != nil { return err }
    tmp := fs.path(p.ID)+".tmp"
    if err := os.WriteFile(tmp, data, 0o644); err != nil { return err }
    return os.Rename(tmp, fs.path(p.ID))
}

// simple id generator
func newID() string {
    b := make([]byte, 8)
    _, _ = rand.Read(b)
    return hex.EncodeToString(b)
}
