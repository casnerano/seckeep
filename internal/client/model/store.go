package model

import (
	"time"

	"github.com/casnerano/seckeep/internal/pkg/model"
)

// StoreData структура записи в локальном хранилище.
type StoreData struct {
	UUID      string         `json:"uuid,omitempty"`
	Type      model.DataType `json:"type"`
	Value     []byte         `json:"value"`
	Version   time.Time      `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	Deleted   bool           `json:"deleted"`
}
