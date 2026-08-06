package upload

import (
	"errors"
	"io"
	"log"
	"sharev2/internal/storage"
	"sync"
	"time"
)

type Session struct {
	mu     sync.Mutex
	upload *storage.Upload

	CreatedAt  time.Time
	LastSeenAt time.Time
}

type Manager struct {
	mu       sync.Mutex // protect sessions
	sessions map[string]*Session
	storage  *storage.Storage

	sessionTimeout time.Duration
}

var (
	ErrNotFound = errors.New("session not found")
)

func CreateManager(storage *storage.Storage) (*Manager, error) {
	sessions := make(map[string]*Session)

	return &Manager{
		sessions:       sessions,
		storage:        storage,
		sessionTimeout: 300 * time.Second,
	}, nil
}

func (m *Manager) CreateSession() (string, error) {
	u, err := m.storage.CreateUpload()
	if err != nil {
		return "", err
	}

	session := &Session{
		upload:     u,
		CreatedAt:  time.Now(),
		LastSeenAt: time.Now(),
	}

	m.mu.Lock()
	m.sessions[u.Id()] = session
	m.mu.Unlock()

	return u.Id(), nil
}

func (m *Manager) Append(id string, offset int64, r io.Reader) (int64, error) {
	session, err := m.findSession(id)
	if err != nil {
		return 0, err
	}

	session.LastSeenAt = time.Now()

	return session.append(offset, r)
}

func (m *Manager) Commit(id string) error {
	session, err := m.findSession(id)
	if err != nil {
		return err
	}

	session.LastSeenAt = time.Now()

	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()

	return session.commit()
}

func (m *Manager) NextOffset(id string) (int64, error) {
	session, err := m.findSession(id)
	if err != nil {
		return 0, err
	}

	return session.upload.Size()
}

func (m *Manager) Cleanup() {
	var expired []*Session

	m.mu.Lock()
	for id, s := range m.sessions {
		if time.Since(s.LastSeenAt) > m.sessionTimeout {
			delete(m.sessions, id)
			expired = append(expired, s)
		}
	}
	m.mu.Unlock()

	for _, s := range expired {
		if err := s.abort(); err != nil {
			log.Printf("cleanup session: %v", err)
		}
	}
}

func (m *Manager) findSession(id string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		return nil, ErrNotFound
	}

	return session, nil
}

func (s *Session) append(offset int64, r io.Reader) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Append(offset, r)
}

func (s *Session) commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Commit()
}

func (s *Session) abort() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Abort()
}
