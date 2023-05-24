// Package handler содержит хендлеры логгера.
package handler

import (
	"bufio"
	"os"

	"github.com/casnerano/seckeep/pkg/log"
)

var _ log.Handler = (*stdOutHandler)(nil)

// Formatter используется для форматеров записи лога.
type Formatter interface {
	Format(record *log.Record) []byte
}

// stdOutHandler обработчик логов для потока stdout.
type stdOutHandler struct {
	formatter Formatter
	level     log.Level
	bubble    bool
}

// NewStdOut конструктор обработчика для потока stdout.
//
// В качестве параметров получает:
// formatter — определяющий формат записи лога,
// level — уровень логирования, который будет учитываться при обработке,
// bubble — всплытие обработчиков (нужно ли применять следующие обработчики).
func NewStdOut(formatter Formatter, level log.Level, bubble bool) log.Handler {
	return &stdOutHandler{formatter, level, bubble}
}

// Handle обрабатывает буфер записей логов.
func (s *stdOutHandler) Handle(records []*log.Record) bool {
	writer := bufio.NewWriter(os.Stdout)

	for _, record := range records {
		if record.Level > s.Level() {
			continue
		}
		_, err := writer.Write(s.formatter.Format(record))
		if err != nil {
			return false
		}
		_ = writer.WriteByte('\n')
	}

	err := writer.Flush()
	return err != nil
}

// Level возвращает уровень логирования.
func (s *stdOutHandler) Level() log.Level {
	return s.level
}

// IsBubble возвращает булево, нужно ли всплытие обработчиков.
func (s *stdOutHandler) IsBubble() bool {
	return s.bubble
}

// Close закрывает ресурс.
func (s *stdOutHandler) Close() error {
	return nil
}
