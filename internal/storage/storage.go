package storage

import (
	"crypto/rand"
	"encoding/hex"
)

type Storage struct {
	RootPath string
}

func generateID() (string, error) {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
