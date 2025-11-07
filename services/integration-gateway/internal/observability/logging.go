package observability

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel nível de log
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger structured logger
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger cria um novo logger
func NewLogger(levelStr string) *Logger {
	var level LogLevel
	switch levelStr {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	default:
		level = InfoLevel
	}

	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug log debug
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= DebugLevel {
		l.log("DEBUG", msg, fields...)
	}
}

// Info log info
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= InfoLevel {
		l.log("INFO", msg, fields...)
	}
}

// Warn log warning
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= WarnLevel {
		l.log("WARN", msg, fields...)
	}
}

// Error log error
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= ErrorLevel {
		l.log("ERROR", msg, fields...)
	}
}

// log método interno de logging estruturado
func (l *Logger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	// Formato: timestamp level message field1=value1 field2=value2
	output := fmt.Sprintf("%s %s %s", timestamp, level, msg)

	// Adiciona fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fields[i]
			value := fields[i+1]
			output += fmt.Sprintf(" %v=%v", key, value)
		}
	}

	l.logger.Println(output)
}

// WithFields cria log com fields
func (l *Logger) WithFields(fields ...interface{}) *LogEntry {
	return &LogEntry{
		logger: l,
		fields: fields,
	}
}

// LogEntry entrada de log com fields
type LogEntry struct {
	logger *Logger
	fields []interface{}
}

// Debug log debug com fields
func (e *LogEntry) Debug(msg string) {
	e.logger.Debug(msg, e.fields...)
}

// Info log info com fields
func (e *LogEntry) Info(msg string) {
	e.logger.Info(msg, e.fields...)
}

// Warn log warn com fields
func (e *LogEntry) Warn(msg string) {
	e.logger.Warn(msg, e.fields...)
}

// Error log error com fields
func (e *LogEntry) Error(msg string) {
	e.logger.Error(msg, e.fields...)
}

// DefaultLogger logger global
var DefaultLogger = NewLogger("info")

// SetLogLevel configura nível de log global
func SetLogLevel(level string) {
	DefaultLogger = NewLogger(level)
}

// Funções de conveniência
func Debug(msg string, fields ...interface{}) {
	DefaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	DefaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	DefaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	DefaultLogger.Error(msg, fields...)
}

func WithFields(fields ...interface{}) *LogEntry {
	return DefaultLogger.WithFields(fields...)
}
