package log

import "time"

type options struct {
	procInterval time.Duration
}

var defaultOptions = options{
	procInterval: time.Millisecond * 100,
}

// Option тип опций для настройки логирования.
type Option func(*options)

// WithProcInterval опция для установки интегрвала запуска обработки буфера логов.
func WithProcInterval(d time.Duration) Option {
	return func(opts *options) {
		opts.procInterval = d
	}
}
