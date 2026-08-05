package storage

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type SaveResult struct {
	ID   string
	Size int64
}

func (s *Storage) SaveStream(src io.Reader) (*SaveResult, error) {
	var (
		id        string
		tmpPath   string
		finalPath string
		file      *os.File
		err       error
	)

	for i := 0; i < 10; i++ {
		id, err = generateID()
		if err != nil {
			return nil, err
		}

		tmpPath = filepath.Join(s.RootPath, id+".uploading")
		finalPath = filepath.Join(s.RootPath, id)

		if _, err := os.Stat(finalPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		if _, err := os.Stat(tmpPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		file, err = os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			return nil, err
		}
		break
	}

	defer file.Close()

	size, err := io.Copy(file, src)

	if err != nil {
		os.Remove(tmpPath)
		return nil, err
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return nil, err
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		return nil, err
	}

	return &SaveResult{
		ID:   id,
		Size: size,
	}, nil
}

func generateID() (string, error) {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
