package log

// Stub структура заглушки логгера.
type Stub struct{}

// NewStub конструктор заглушки логгера.
func NewStub() *Stub {
	return &Stub{}
}

// Log stub-метод.
func (s Stub) Log(level Level, message string, context ...any) {}

// Emergency stub-метод.
func (s Stub) Emergency(message string, context ...any) {}

// Alert stub-метод.
func (s Stub) Alert(message string, context ...any) {}

// Critical stub-метод.
func (s Stub) Critical(message string, context ...any) {}

// Error stub-метод.
func (s Stub) Error(message string, context ...any) {}

// Warning stub-метод.
func (s Stub) Warning(message string, context ...any) {}

// Notice stub-метод.
func (s Stub) Notice(message string, context ...any) {}

// Info stub-метод.
func (s Stub) Info(message string, context ...any) {}

// Debug stub-метод.
func (s Stub) Debug(message string, context ...any) {}
