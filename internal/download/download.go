package download

import (
	"encoding/hex"
	"errors"
	"log"
	"os"
	"sharev2/internal/model"
	"sharev2/internal/storage"
	"time"
)

type Manager struct {
	storage *storage.Storage
}

var (
	ErrNotFound     = errors.New("not found")
	ErrMetaDataMiss = errors.New("metadata miss")
)

const (
	expiration = 7 * 24 * time.Hour
)

func NewManager(storage *storage.Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

func (m *Manager) OpenDownload(clientId string) (*os.File, *model.MetaData, error) {
	if len(clientId) != 32 {
		return nil, nil, ErrNotFound
	}

	buf, err := hex.DecodeString(clientId)
	if err != nil {
		return nil, nil, ErrNotFound
	}

	id := hex.EncodeToString(buf)

	file, err := m.storage.OpenData(id)
	if err != nil {
		return nil, nil, ErrNotFound
	}

	metaData, err := m.storage.ReadMetaData(id)
	if err != nil {
		file.Close()
		return nil, nil, ErrMetaDataMiss
	}

	metaData.LastDownload = time.Now()

	if err := m.storage.WriteMetaData(&metaData); err != nil {
		log.Println(err)
	}

	return file, &metaData, nil
}

func (m *Manager) AddDownloadCount(clientId string) error {
	if len(clientId) != 32 {
		return ErrNotFound
	}

	buf, err := hex.DecodeString(clientId)
	if err != nil {
		return ErrNotFound
	}

	id := hex.EncodeToString(buf)

	metaData, err := m.storage.ReadMetaData(id)
	metaData.DownloadCount++

	return m.storage.WriteMetaData(&metaData)
}

func (m *Manager) Cleanup() {
	metaData, err := m.ListMetaData()
	if err != nil {
		log.Println(err)
	}

	for _, data := range metaData {
		if time.Since(data.LastDownload) <= expiration {
			continue
		}

		if err := m.storage.Delete(data.ID); err != nil {
			log.Println(err)
		}
	}
}
