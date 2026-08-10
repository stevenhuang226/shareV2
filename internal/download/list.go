package download

import "sharev2/internal/model"

func (m *Manager) ListMetaData() ([]model.MetaData, error) {
	return m.storage.ListMetaData()
}
