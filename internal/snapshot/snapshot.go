package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Snapshot holds a point-in-time record of observed port bindings.
type Snapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Ports     []portscanner.PortInfo `json:"ports"`
}

// Store persists and retrieves snapshots from disk.
type Store struct {
	path string
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk as JSON, overwriting any existing file.
func (s *Store) Save(snap Snapshot) error {
	f, err := os.CreateTemp("", "portwatch-snap-*.json")
	if err != nil {
		return err
	}
	tmpName := f.Name()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	f.Close()

	return os.Rename(tmpName, s.path)
}

// Load reads the most recently saved snapshot from disk.
// Returns an empty Snapshot and no error when the file does not yet exist.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Diff compares two snapshots and returns ports that were added or removed.
func Diff(prev, curr Snapshot) (added, removed []portscanner.PortInfo) {
	prevMap := make(map[string]portscanner.PortInfo, len(prev.Ports))
	for _, p := range prev.Ports {
		prevMap[p.Address] = p
	}

	currMap := make(map[string]portscanner.PortInfo, len(curr.Ports))
	for _, p := range curr.Ports {
		currMap[p.Address] = p
	}

	for addr, p := range currMap {
		if _, exists := prevMap[addr]; !exists {
			added = append(added, p)
		}
	}
	for addr, p := range prevMap {
		if _, exists := currMap[addr]; !exists {
			removed = append(removed, p)
		}
	}
	return added, removed
}
