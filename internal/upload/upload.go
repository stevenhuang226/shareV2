package upload

import (
	"io"
	"sharev2/internal/storage"
	"sync"
)

type Session struct {
	mu         sync.Mutex
	upload     *storage.Upload
	expectSize int64
}

type Manager struct {
	mu       sync.Mutex // protect sessions
	sessions map[string]*Session
}

func (m *Manager) NewSession() (string, error) {
	/*
		create new session

		write with sync lock
	*/
}

func (s *Session) Append(offset int64, r io.Reader) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Append(offset, r)
}

func (s *Session) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Commit()
}

func (s *Session) Abort() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.upload.Abort()
}
