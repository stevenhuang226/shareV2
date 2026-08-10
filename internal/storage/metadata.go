package storage

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sharev2/internal/model"
)

func (s *Storage) WriteMetaData(metaData *model.MetaData) error {
	metaPath := filepath.Join(s.MetaDataRoot, metaData.ID+".json")

	file, err := os.Create(metaPath)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(metaData)
}

func (s *Storage) ReadMetaData(id string) (model.MetaData, error) {
	metaPath := filepath.Join(s.MetaDataRoot, id+".json")

	metaData := model.MetaData{
		ID:       id,
		Name:     "ERROR",
		MIMEType: "ERROR",
	}

	file, err := os.Open(metaPath)
	if err != nil {
		return metaData, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&metaData); err != nil {
		return metaData, err
	}

	return metaData, nil
}

func (s *Storage) ListMetaData() ([]model.MetaData, error) {
	entries, err := os.ReadDir(s.MetaDataRoot)
	if err != nil {
		return nil, err
	}

	metaData := []model.MetaData{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		id, ok := metaDataId(entry.Name())
		if !ok {
			continue
		}

		data, err := s.ReadMetaData(id)
		if err != nil {
			return nil, err
		}

		metaData = append(metaData, data)
	}

	return metaData, nil
}

func metaDataId(name string) (string, bool) {
	if len(name) != 37 || name[32:] != ".json" {
		return "", false
	}

	id := name[:32]

	if _, err := hex.DecodeString(id); err != nil {
		return "", false
	}

	return id, true
}
