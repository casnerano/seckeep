package model

import (
	"github.com/casnerano/seckeep/internal/pkg/model"
)

// DataTypeable интерфейс секретных данных.
type DataTypeable interface {
	Type() model.DataType
}

// DataCredential структура учетной записи.
type DataCredential struct {
	Login    string   `json:"login" validate:"required"`
	Password string   `json:"password" validate:"required"`
	Meta     []string `json:"meta"`
}

// Type возвращает тип структуры.
func (c DataCredential) Type() model.DataType {
	return model.DataTypeCredential
}

// DataText структура простого текста.
type DataText struct {
	Value string   `json:"value" validate:"required"`
	Meta  []string `json:"meta"`
}

// Type возвращает тип структуры.
func (c DataText) Type() model.DataType {
	return model.DataTypeText
}

// DataCard структура банковской карты.
type DataCard struct {
	Number    string   `json:"number" validate:"required,credit_card"`
	MonthYear string   `json:"month_year" validate:"required,datetime=01.02"`
	CVV       string   `json:"cvv" validate:"required"`
	Owner     string   `json:"owner"`
	Meta      []string `json:"meta"`
}

// Type возвращает тип структуры.
func (c DataCard) Type() model.DataType {
	return model.DataTypeCard
}

// DataDocument структура документа.
type DataDocument struct {
	Name    string   `json:"name" validate:"required"`
	Content []byte   `json:"content" validate:"required"`
	Meta    []string `json:"meta"`
}

// Type возвращает тип структуры.
func (c DataDocument) Type() model.DataType {
	return model.DataTypeDocument
}
