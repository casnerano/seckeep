package log

// Level тип уровеня лога.
type Level int

// Варианты уровней логирования.
const (
	LevelEmergency Level = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
)

var levelString = map[Level]string{
	LevelEmergency: "EMERGENCY",
	LevelAlert:     "ALERT",
	LevelCritical:  "CRITICAL",
	LevelError:     "ERROR",
	LevelWarning:   "WARNING",
	LevelNotice:    "NOTICE",
	LevelInfo:      "INFO",
	LevelDebug:     "DEBUG",
}

// String метод выводит уровень логирования в виде строки.
func (l Level) String() string {
	return levelString[l]
}
