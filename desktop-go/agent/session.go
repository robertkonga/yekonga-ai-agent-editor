package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Session struct {
	ID          string        `json:"id"`
	Provider    string        `json:"provider"`
	History     []ChatMessage `json:"history"`
	LastUpdated time.Time     `json:"last_updated"`
}

type SessionManager struct {
	storageDir string
	sessions   map[string]*Session
	mu         sync.RWMutex
}

func NewSessionManager(storageDir string) (*SessionManager, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}
	return &SessionManager{
		storageDir: storageDir,
		sessions:   make(map[string]*Session),
	}, nil
}

func (m *SessionManager) GetSession(id string) (*Session, error) {
	m.mu.RLock()
	if s, ok := m.sessions[id]; ok {
		m.mu.RUnlock()
		return s, nil
	}
	m.mu.RUnlock()

	// Try to load from disk
	path := filepath.Join(m.storageDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not found
		}
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.sessions[id] = &s
	m.mu.Unlock()

	return &s, nil
}

func (m *SessionManager) SaveSession(s *Session) error {
	s.LastUpdated = time.Now()
	
	m.mu.Lock()
	m.sessions[s.ID] = s
	m.mu.Unlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(m.storageDir, s.ID+".json")
	return os.WriteFile(path, data, 0644)
}

func (m *SessionManager) ListSessions() ([]*Session, error) {
	entries, err := os.ReadDir(m.storageDir)
	if err != nil {
		return nil, err
	}

	var sessions []*Session
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			id := strings.TrimSuffix(e.Name(), ".json")
			s, err := m.GetSession(id)
			if err == nil && s != nil {
				sessions = append(sessions, s)
			}
		}
	}
	return sessions, nil
}
