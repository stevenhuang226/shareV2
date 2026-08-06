package upload

import (
	"errors"
	"io"
	"sharev2/internal/storage"
	"sync"
)

type Session struct {
	mu     sync.Mutex
	upload *storage.Upload
}

type Manager struct {
	mu       sync.Mutex // protect sessions
	sessions map[string]*Session
	storage  *storage.Storage
}

func CreateManager(storage *storage.Storage) (*Manager, error) {
	sessions := make(map[string]*Session)

	return &Manager{
		sessions: sessions,
		storage:  storage,
	}, nil
}

func (m *Manager) NewSession() (string, error) {
	u, err := m.storage.CreateUpload()
	if err != nil {
		return "", err
	}

	session := &Session{
		upload: u,
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

	return session.append(offset, r)
}

func (m *Manager) Commit(id string) error {
	session, err := m.findSession(id)
	if err != nil {
		return err
	}

	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()

	return session.commit()
}

func (m *Manager) findSession(id string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
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
