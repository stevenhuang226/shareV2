package storage

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
)

type Storage struct {
	DataRoot     string
	MetaDataRoot string
}

func generateID() (string, error) {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func (s *Storage) Delete(id string) error {
	dataPath := filepath.Join(s.DataRoot, id)
	metaDataPath := filepath.Join(s.MetaDataRoot, id+".json")

	if err := os.Remove(metaDataPath); err != nil {
		return err
	}

	if err := os.Remove(dataPath); err != nil {
		return err
	}

	return nil
}
