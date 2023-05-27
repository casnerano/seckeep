package model

import (
	"time"

	"github.com/casnerano/seckeep/internal/pkg"
)

// DataCreateRequest структура запроса создания данных.
type DataCreateRequest struct {
	Type      pkg.DataType `json:"type" validate:"required,enum"`
	Value     []byte       `json:"value" validate:"required"`
	Version   time.Time    `json:"version" validate:"required"`
	CreatedAt time.Time    `json:"created_at" validate:"required"`
}

// DataUpdateRequest структура запроса обновления данных.
type DataUpdateRequest struct {
	Value   []byte    `json:"value" validate:"required"`
	Version time.Time `json:"version" validate:"required"`
}
