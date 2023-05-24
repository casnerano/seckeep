package formatter

import (
	"fmt"

	"github.com/casnerano/seckeep/pkg/log"
	"github.com/casnerano/seckeep/pkg/log/handler"
)

const layoutDate = "02.01.2006 15:04:05"

type fText struct{}

// NewText конструктор текстового форматера.
func NewText() handler.Formatter {
	return &fText{}
}

// Format возвращает слайс байт текстового представления лога.
func (j fText) Format(record *log.Record) []byte {
	text := fmt.Sprintf(
		"[%s][%s] %s %s",
		record.Level.String(),
		record.Channel,
		record.Date.Format(layoutDate),
		record.Message,
	)

	if len(record.Context) > 0 {
		text = fmt.Sprintf("%s %v", text, record.Context)
	}

	return []byte(text)
}
