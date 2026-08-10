package storage

import (
	"os"
	"path/filepath"
)

func (s *Storage) OpenData(id string) (*os.File, error) {
	path := filepath.Join(s.DataRoot, id)

	return os.Open(path)
}
