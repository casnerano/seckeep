// Package model содержит общие модели данных.
package model

import "time"

// Data структура секретных данных.
type Data struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	Type      DataType  `json:"type"`
	Value     []byte    `json:"value"`
	Version   time.Time `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// DataType типы данных.
type DataType string

// IsValid проверяет на валидность тип данных.
func (d DataType) IsValid() bool {
	switch d {
	case DataTypeCredential, DataTypeText,
		DataTypeCard, DataTypeDocument:
		return true
	}
	return false
}

// Варианты типов данных.
const (
	// DataTypeCredential учетная запись.
	DataTypeCredential DataType = "CREDENTIAL"

	// DataTypeText простой текст.
	DataTypeText DataType = "TEXT"

	// DataTypeCard банкоская карта.
	DataTypeCard DataType = "CARD"

	// DataTypeDocument произвольный документ.
	DataTypeDocument DataType = "DOCUMENT"
)
