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
