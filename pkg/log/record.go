package log

import (
	"time"
)

// Record структура записи лога.
type Record struct {
	Level   Level     `json:"level"`
	Date    time.Time `json:"date"`
	Channel string    `json:"channel"`
	Message string    `json:"message"`
	Context []any     `json:"context,omitempty"`
}
