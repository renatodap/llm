package conversation

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// FileStore persists threads as JSON files under a directory
type FileStore struct {
    dir string
    mu  sync.Mutex
}

func NewFileStore(dir string) (*FileStore, error) {
    if dir == "" { dir = ".llm_threads" }
    if err := os.MkdirAll(dir, 0o755); err != nil { return nil, err }
    return &FileStore{dir: dir}, nil
}

func (fs *FileStore) path(id string) string { return filepath.Join(fs.dir, id+".json") }

func (fs *FileStore) CreateThread(title string) (*Thread, error) {
    fs.mu.Lock(); defer fs.mu.Unlock()
    id := newID()
    now := time.Now()
    t := &Thread{
        ID: id,
        Title: title,
        CreatedAt: now,
        UpdatedAt: now,
        Metadata: make(map[string]interface{}),
        Messages: []Message{},
    }
    if err := fs.save(t); err != nil { return nil, err }
    return t, nil
}

// CreateThreadForProject creates a thread scoped to a project
func (fs *FileStore) CreateThreadForProject(projectID, title string) (*Thread, error) {
    fs.mu.Lock(); defer fs.mu.Unlock()
    id := newID()
    now := time.Now()
    t := &Thread{
        ID: id,
        ProjectID: projectID,
        Title: title,
        CreatedAt: now,
        UpdatedAt: now,
        Metadata: make(map[string]interface{}),
        Messages: []Message{},
    }
    if err := fs.save(t); err != nil { return nil, err }
    return t, nil
}

func (fs *FileStore) GetThread(id string) (*Thread, error) {
    data, err := os.ReadFile(fs.path(id))
    if err != nil { return nil, err }
    var t Thread
    if err := json.Unmarshal(data, &t); err != nil { return nil, err }
    return &t, nil
}

func (fs *FileStore) ListThreads() ([]*Thread, error) {
    entries, err := os.ReadDir(fs.dir)
    if err != nil { return nil, err }
    var out []*Thread
    for _, e := range entries {
        if e.IsDir() { continue }
        if filepath.Ext(e.Name()) != ".json" { continue }
        data, err := os.ReadFile(filepath.Join(fs.dir, e.Name()))
        if err != nil { continue }
        var t Thread
        if json.Unmarshal(data, &t) == nil { out = append(out, &t) }
    }
    return out, nil
}

func (fs *FileStore) AppendMessage(threadID string, msg Message) (*Thread, error) {
    fs.mu.Lock(); defer fs.mu.Unlock()
    t, err := fs.GetThread(threadID)
    if err != nil { return nil, err }
    if msg.ID == "" { msg.ID = newID() }
    if msg.CreatedAt.IsZero() { msg.CreatedAt = time.Now() }
    t.Messages = append(t.Messages, msg)
    t.UpdatedAt = time.Now()
    if err := fs.save(t); err != nil { return nil, err }
    return t, nil
}

func (fs *FileStore) UpdateSummary(threadID string, summary string) error {
    fs.mu.Lock(); defer fs.mu.Unlock()
    t, err := fs.GetThread(threadID)
    if err != nil { return err }
    t.Summary = summary
    t.UpdatedAt = time.Now()
    return fs.save(t)
}

// UpdateThread saves the entire thread (used for metadata updates)
func (fs *FileStore) UpdateThread(t *Thread) error {
    fs.mu.Lock(); defer fs.mu.Unlock()
    t.UpdatedAt = time.Now()
    return fs.save(t)
}

func (fs *FileStore) save(t *Thread) error {
    data, err := json.MarshalIndent(t, "", "  ")
    if err != nil { return err }
    tmp := fs.path(t.ID)+".tmp"
    if err := os.WriteFile(tmp, data, 0o644); err != nil { return err }
    return os.Rename(tmp, fs.path(t.ID))
}

func newID() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil { return fmt.Sprintf("%d", time.Now().UnixNano()) }
    return hex.EncodeToString(b)
}
