package storage

import (
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
