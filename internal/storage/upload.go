package storage

import (
	"errors"
	"os"
	"path/filepath"
)

type Upload struct {
	file      *os.File
	id        string
	size      int64
	tmpPath   string
	finalPath string
}

const (
	namingRetryLimit int = 32
)

func (s *Storage) CreateUpload() (*Upload, error) {
	var (
		id   string
		err  error
		file *os.File
	)
	for i := 0; i < namingRetryLimit; i++ {
		id, err = generateID()
		if err != nil {
			return nil, err
		}

		tmp := filepath.Join(s.RootPath, id+".uploading")
		final := filepath.Join(s.RootPath, id)

		if _, err := os.Stat(tmp); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		if _, err := os.Stat(final); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		file, err = os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			return nil, err
		}

		return &Upload{
			file:      file,
			id:        id,
			size:      0,
			tmpPath:   tmp,
			finalPath: final,
		}, nil
	}

	return nil, errors.New("naming retry limit")
}
