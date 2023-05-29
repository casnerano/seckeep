package formatter

import (
	"encoding/json"

	"github.com/casnerano/seckeep/pkg/log"
	"github.com/casnerano/seckeep/pkg/log/handler"
)

type fJSON struct{}

// NewJSON конструктор JSON форматера.
func NewJSON() handler.Formatter {
	return &fJSON{}
}

// Format возвращает слайс байт JSON представления лога.
func (j fJSON) Format(record *log.Record) []byte {
	b, _ := json.Marshal(record)
	return b
}
