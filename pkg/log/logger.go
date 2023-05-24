package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Handler используется для обработчиков логирования.
//
// Задача обработчика заключается в том, чтобы обеспечить запись буфера логов,
// при необходимости, установить поведения всплытия для последующих обрабочиков, и закрытие открытых ресурсов.
type Handler interface {
	Handle(records []*Record) bool

	Level() Level
	IsBubble() bool

	Close() error
}

// Logger структура логгера.
type Logger struct {
	channel string

	mu       sync.Mutex
	records  []*Record
	handlers []Handler

	done chan struct{}
	stop chan struct{}
	opts options
}

// Loggable интерфейс логгера.
type Loggable interface {
	Log(level Level, message string, context ...any)
	Emergency(message string, context ...any)
	Alert(message string, context ...any)
	Critical(message string, context ...any)
	Error(message string, context ...any)
	Warning(message string, context ...any)
	Notice(message string, context ...any)
	Info(message string, context ...any)
	Debug(message string, context ...any)
}

// New конструктор структуры логирования.
//
// Позволяет устанавливать произвольные опции — opts, например:
//
//	New(WithProcInterval(time.Millisecond * 100))
func New(channel string, opts ...Option) *Logger {
	logOpts := defaultOptions
	for _, opt := range opts {
		opt(&logOpts)
	}

	l := &Logger{
		channel: channel,
		records: make([]*Record, 0),
		done:    make(chan struct{}),
		stop:    make(chan struct{}),
		opts:    logOpts,
	}

	l.processing()
	return l
}

// processing запускает обработку буфера логов по заданному интервалу.
func (l *Logger) processing() {
	go func() {
		defer close(l.done)
		for {
			select {
			case <-time.After(l.opts.procInterval):
				l.flush()
			case <-l.stop:
				l.flush()
				return
			}
		}
	}()
}

// flush отдает буфер логов хендлерам для обработки, и очищает.
func (l *Logger) flush() {
	l.mu.Lock()
	records := l.records
	l.records = make([]*Record, 0, len(records))
	l.mu.Unlock()

	for _, handler := range l.handlers {
		handler.Handle(records)
		if !handler.IsBubble() {
			break
		}
	}
}

// AddHandler добавляет обработчик в список обрабочиков для логирования.
func (l *Logger) AddHandler(handler Handler) {
	l.handlers = append(l.handlers, handler)
}

// Close изящно останавливает горутину обработки буфера логов,
// и закрывает ресуры обработчиков.
func (l *Logger) Close() error {
	close(l.stop)
	<-l.done
	var err error
	for k := range l.handlers {
		err = l.handlers[k].Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close %T handler\n", l.handlers[k])
		}
	}
	return nil
}

// Log добавляет запись с заданным уровнем в бефер логов.
func (l *Logger) Log(level Level, message string, context ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, &Record{
		level,
		time.Now(),
		l.channel,
		message,
		context,
	})
}

// Emergency логирование с уровнем Emergency.
func (l *Logger) Emergency(message string, context ...any) {
	l.Log(LevelEmergency, message, context...)
}

// Alert логирование с уровнем Alert.
func (l *Logger) Alert(message string, context ...any) {
	l.Log(LevelAlert, message, context...)
}

// Critical логирование с уровнем Critical.
func (l *Logger) Critical(message string, context ...any) {
	l.Log(LevelCritical, message, context...)
}

// Error логирование с уровнем Error.
func (l *Logger) Error(message string, context ...any) {
	l.Log(LevelError, message, context...)
}

// Warning логирование с уровнем Warning.
func (l *Logger) Warning(message string, context ...any) {
	l.Log(LevelWarning, message, context...)
}

// Notice логирование с уровнем Notice.
func (l *Logger) Notice(message string, context ...any) {
	l.Log(LevelNotice, message, context...)
}

// Info логирование с уровнем Info.
func (l *Logger) Info(message string, context ...any) {
	l.Log(LevelInfo, message, context...)
}

// Debug логирование с уровнем Debug.
func (l *Logger) Debug(message string, context ...any) {
	l.Log(LevelDebug, message, context...)
}
