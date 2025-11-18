package conversation

import (
    "sort"
)

// VectorIndex is a simple in-memory store of embeddings for messages
type VectorIndex struct {
    vecs map[string][]float64 // messageID -> vector
}

func NewVectorIndex() *VectorIndex { return &VectorIndex{vecs: make(map[string][]float64)} }

func (vi *VectorIndex) Add(messageID string, embedding []float64) { vi.vecs[messageID] = embedding }

// Query returns top-k messageIDs by cosine similarity
func (vi *VectorIndex) Query(query []float64, k int, sim func(a, b []float64) float64) []string {
    type rec struct{ id string; s float64 }
    tmp := make([]rec, 0, len(vi.vecs))
    for id, v := range vi.vecs {
        tmp = append(tmp, rec{id: id, s: sim(query, v)})
    }
    sort.Slice(tmp, func(i, j int) bool { return tmp[i].s > tmp[j].s })
    if k > len(tmp) { k = len(tmp) }
    out := make([]string, 0, k)
    for i := 0; i < k; i++ { out = append(out, tmp[i].id) }
    return out
}

