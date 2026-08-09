package storage

import (
	"errors"
	"io"
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

var (
	ErrFileIsNil      = errors.New("file is nil")
	ErrOffsetWrong    = errors.New("wrong offset")
	ErrOffsetMismatch = errors.New("append only")
	ErrFileNotClose   = errors.New("file not close")
)

const (
	namingRetryLimit int   = 32
	maxFileSize      int64 = 10 * 1024 * 1024 * 1024
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

		tmp := filepath.Join(s.DataRoot, id+".uploading")
		final := filepath.Join(s.DataRoot, id)

		if _, err := os.Stat(tmp); err == nil {
			continue // tmp already exist, retry
		} else if !os.IsNotExist(err) {
			return nil, err // os.Stat err
		}

		if _, err := os.Stat(final); err == nil {
			continue // final already exist, retry
		} else if !os.IsNotExist(err) {
			return nil, err // os.Stat err
		}

		file, err = os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				continue // already exist, retry
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

func (u *Upload) Id() string {
	return u.id
}

func (u *Upload) Size() int64 {
	return u.size
}

/* write bytes int64, Err error */
func (u *Upload) Append(offset int64, r io.Reader) (int64, error) {
	if u.file == nil {
		return 0, ErrFileIsNil
	}
	if offset < 0 || offset > maxFileSize {
		return 0, ErrOffsetWrong
	}

	/* not support offset yet */
	if offset != u.size {
		return 0, ErrOffsetMismatch
	}

	remaining := maxFileSize - offset // bytes it can use

	lr := &io.LimitedReader{
		R: r,
		N: remaining,
	} // limited to read "remaining" bytes data

	size, err := io.Copy(u.file, lr)
	u.size += size

	if err != nil {
		return size, err // io.Copy err
	}

	if lr.N <= 0 {
		return size, errors.New("limit size") // lr.N <= 0 (hit size limit)
	}

	return size, nil
}

func (u *Upload) Commit() error {
	if u.file == nil {
		return ErrFileIsNil
	}

	if err := u.file.Sync(); err != nil {
		return err
	}

	if err := u.file.Close(); err != nil {
		return err
	}

	u.file = nil

	if err := os.Rename(u.tmpPath, u.finalPath); err != nil {
		return err
	}

	return nil
}

func (u *Upload) RenameTmp() error {
	if u.file != nil {
		return ErrFileNotClose
	}

	if err := os.Rename(u.tmpPath, u.finalPath); err != nil {
		return err
	}

	return nil
}

func (u *Upload) Abort() error {
	if u.file == nil {
		return ErrFileIsNil
	}

	if err := u.file.Close(); err != nil {
		return err
	}

	u.file = nil

	if err := os.Remove(u.tmpPath); err != nil {
		return err
	}

	return nil
}
