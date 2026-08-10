package download

import "sharev2/internal/model"

func (m *Manager) ListMetaDatas() ([]model.MetaData, error) {
	return m.storage.ListMetaData()
}
